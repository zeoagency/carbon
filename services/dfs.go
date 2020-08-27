package services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/zeoagency/carbon/helpers"
	"github.com/zeoagency/carbon/models"
)

// countries keeps country codes for DataForSEO.
var countries = `{"AF":"2004","AL":"2008","DZ":"2012","AS":"2016","AD":"2020","AO":"2024","AG":"2028","AZ":"2031","AR":"2032","AU":"2036","AT":"2040","BS":"2044","BH":"2048","BD":"2050","AM":"2051","BB":"2052","BE":"2056","BM":"2060","BT":"2064","BO":"2068","BA":"2070","BW":"2072","BR":"2076","BZ":"2084","SB":"2090","VG":"2092","BN":"2096","BG":"2100","MM":"2104","BI":"2108","BY":"2112","KH":"2116","CM":"2120","CA":"2124","CV":"2132","KY":"2136","CF":"2140","LK":"2144","TD":"2148","CL":"2152","CN":"2156","TW":"2158","CO":"2170","CG":"2178","CD":"2180","CK":"2184","CR":"2188","HR":"2191","CY":"2196","CZ":"2203","BJ":"2204","DK":"2208","DM":"2212","DO":"2214","EC":"2218","SV":"2222","ET":"2231","EE":"2233","FJ":"2242","FI":"2246","FR":"2250","DJ":"2262","GA":"2266","GE":"2268","GM":"2270","PS":"2275","DE":"2276","GH":"2288","GI":"2292","KI":"2296","GR":"2300","GL":"2304","GP":"2312","GU":"2316","GT":"2320","GY":"2328","HT":"2332","HN":"2340","HK":"2344","HU":"2348","IS":"2352","IN":"2356","ID":"2360","IQ":"2368","IE":"2372","IL":"2376","IT":"2380","CI":"2384","JM":"2388","JP":"2392","KZ":"2398","JO":"2400","KE":"2404","KR":"2410","KW":"2414","KG":"2417","LA":"2418","LB":"2422","LS":"2426","LV":"2428","LY":"2434","LI":"2438","LT":"2440","LU":"2442","MO":"2446","MG":"2450","MW":"2454","MY":"2458","MV":"2462","ML":"2466","MT":"2470","MU":"2480","MX":"2484","MN":"2496","MD":"2498","ME":"2499","MS":"2500","MA":"2504","MZ":"2508","OM":"2512","NA":"2516","NR":"2520","NP":"2524","NL":"2528","VU":"2548","NZ":"2554","NI":"2558","NE":"2562","NG":"2566","NU":"2570","NF":"2574","NO":"2578","FM":"2583","PK":"2586","PA":"2591","PG":"2598","PY":"2600","PE":"2604","PH":"2608","PN":"2612","PL":"2616","PT":"2620","TL":"2626","PR":"2630","QA":"2634","RO":"2642","RU":"2643","RW":"2646","SH":"2654","AI":"2660","VC":"2670","SM":"2674","ST":"2678","SA":"2682","SN":"2686","RS":"2688","SC":"2690","SL":"2694","SG":"2702","SK":"2703","VN":"2704","SI":"2705","SO":"2706","ZA":"2710","ZW":"2716","ES":"2724","SE":"2752","CH":"2756","TJ":"2762","TH":"2764","TG":"2768","TK":"2772","TO":"2776","TT":"2780","AE":"2784","TN":"2788","TR":"2792","TM":"2795","UG":"2800","UA":"2804","MK":"2807","EG":"2818","GB":"2826","GG":"2831","JE":"2832","TZ":"2834","US":"2840","VI":"2850","BF":"2854","UY":"2858","UZ":"2860","VE":"2862","WS":"2882","ZM":"2894"}`

type dfsApiResponse struct {
	StatusCode int `json:"status_code"`
	Tasks      []struct {
		Result []struct {
			Keyword string `json:"keyword"`
			Items   []struct {
				Type        string `json:"type"`
				Title       string `json:"title"`
				Description string `json:"description"`
				URL         string `json:"url"`
			} `json:"items"`
		} `json:"result"`
	} `json:"tasks"`
}

type dfsApiRequest struct {
	Keyword   string `json:"keyword"`
	Gl        string `json:"location_code"`
	Hl        string `json:"language_code"`
	SerpLimit int    `json:"depth"`
	Device    string `json:"device"`
}

// getResultFromDFSApi returns DFS API response for the given data.
func getResultFromDFSApi(kws keywords, country, language string, serpLimit int) (map[string]*dfsApiResponse, int, error) {
	// set country code.
	c := make(map[string]string)
	_ = json.Unmarshal([]byte(countries), &c)
	cCode := c[strings.ToUpper(country)]
	// a map to keep result
	responses := make(map[string]*dfsApiResponse) // the key is originalURL

	wg := new(sync.WaitGroup)
	for _, kw := range kws.ToStringSlice() {
		rq := []dfsApiRequest{}
		rq = append(rq, dfsApiRequest{
			Keyword:   kw,
			Gl:        cCode,
			Hl:        language,
			SerpLimit: serpLimit,
			Device:    "desktop",
		})

		wg.Add(1)
		go sendRequest(wg, responses, kws.Original(kw), rq)
	}

	wg.Wait()
	return responses, http.StatusCreated, nil
}

// sendRequest works as async, adds all results to the given slice.
func sendRequest(wg *sync.WaitGroup, responses map[string]*dfsApiResponse, key string, rq []dfsApiRequest) {
	defer wg.Done()

	rqJson, err := json.Marshal(rq)
	if err != nil {
		return
	}

	// Create the request.
	req, err := http.NewRequest("POST", "https://api.dataforseo.com/v3/serp/google/organic/live/regular", bytes.NewReader(rqJson))
	if err != nil {
		return
	}

	req.SetBasicAuth(os.Getenv("DFS_API_USER"), os.Getenv("DFS_API_PASSWORD"))

	// Send the request.
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	// Check the result's status code.
	if !(res.StatusCode >= 200 && res.StatusCode <= 299) {
		log.Printf("Error: Unavailable DFS API Service. Status: %d\n", res.StatusCode)
		return
	}

	// Read the result, unmarshal it to the struct.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	response := dfsApiResponse{}
	err = json.Unmarshal(body, &response) // TODO: handle this error.

	// Check the result's (body) status code.
	if !(response.StatusCode >= 20000 && response.StatusCode <= 29999) {
		log.Printf("Error: Unavailable DFS API Service. Status: %d\n", response.StatusCode)
	}

	responses[key] = &response
}

// parseDFSResponseToFieldsForURLs extract the response to the URLSet.
// It only adds to success list when domains are matched.
// If it couldn't find any related URLs, it adds to the fail list.
func parseDFSResponseToFieldsForURLs(responses map[string]*dfsApiResponse, urlSet *models.URLSet) {
	for originalURL, response := range responses {
		for _, task := range response.Tasks {
			r := []string{}
			for _, result := range task.Result {
				for _, item := range result.Items {
					if len(r) == 3 {
						break
					}

					urlDomain, _, err := helpers.ExtractURL(item.URL)
					if err != nil {
						continue // The result's URL is not a valid URL.
					}
					if urlSet.URLs[result.Keyword].BaseURL == urlDomain {
						r = append(r, item.URL)
					}
				}
			}

			if len(r) != 0 {
				urlSet.AddSuccess(originalURL, r)
			} else {
				urlSet.AddFail(originalURL, "We could not find any related URLs.")
			}
		}
	}
}

// parseDFSResponseToFieldsForKeywords extract the response to the KeywordSet.
// It only adds to success list when the value is valid.
// If it couldn't find any results, it adds to the fail list.
func parseDFSResponseToFieldsForKeywords(responses map[string]*dfsApiResponse, keywordSet *models.KeywordSet) {
	for keyword, response := range responses {
		for _, task := range response.Tasks {
			r := []models.KeywordSuccessResult{}
			for _, result := range task.Result {
				for _, item := range result.Items {
					if len(r) == 10 {
						break
					}
					r = append(r, models.KeywordSuccessResult{
						Title: item.Title,
						Desc:  item.Description,
						URL:   item.URL,
					})
				}
			}
			if len(r) != 0 {
				keywordSet.AddSuccess(keyword, r)
			} else {
				keywordSet.AddFail(keyword, "We could not find any related URLs.")
			}
		}
	}
}

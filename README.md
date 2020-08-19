# ZEO.ORG - Carbon API

Carbon aims to find related URLs of given URLs or Keywords.  
It's mostly used to find alternative for 404 pages or SERP operations.  
It exports data in Excel or Google Sheets.

The API is served at AWS Lambda.

## Features

- Supports URL and Keyword options.
	- URL option; finds related 3 URLs.  
	  Mostly used to find alternative URLs for 404 pages.  
	- Keyword option; finds related 10 URLs.  
	  The result includes title, url, and description for each input keywords.  
	  Mostly used for SERP.  
- Supports country and language specification.  
- Supports 2 export options; Excel and Google Sheets.  
	- Suggested URLs are made bold for URL options.  
- Supports internal accounts with limitation.
	- For non-login users, the limit is 100 URLs.

## Endpoint

**URL:** /

**Request:**

- Type: **POST**
- Params: 
	- **type** `must`  
	  options: `keyword` or `url`.  
	  note: `keyword` option is only available for internal users.
	- **format** `must`  
	  options: `excel` or `sheet`.
	- **country** `must`  
	  options: all countries supported by Google. 
	- **langauge** `must`  
	  options: all languages supported by Google.
	- **accountName**  
	- **accountPassword**  
- Header:
	- **Accept**  `must`  
	  If the format is `excel`,  
	  you need to set `Accept: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`  
- Body:
	- Raw Data  
		- As a JSON value,  
		  For example;
			```json
			{
			    "values": [
			      {
			        "value": "https://tools.zeo.org/carbon"
			      },
			      {
			        "value": "https://seo.do/pricing/"
			      },
			      ...
			    ]
			}
			```

**Response:**

Status;

- Type: **201**
	- That means the data is created.
- Type: **400**
	- That means the inputs are not valid.  
	  Check the error message.
- Type: **401**
	- That means auth is not successful.
- Type: **403**
	- That means the method is forbidden.  
	  Use POST method.
- Type: **500**
	- That means internal error occurs while creating the data.
- Type: **503**
	- That means the service is not available.  
	  Try later.

Header and body;

- For **excel**;
	- Header  
		```
		Content-Disposition: attachment; filename="result.xlsx"
		Content-Length: 453646
		Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
		```
	- Body  
		`file`
- For **sheet**;  
	- Body  
		```
		 {
		     "sheetURL": "https://docs.google.com/spreadsheets/d/...",
		 }
		```

## Development

#### Requirements

- SERP API Credantials. (to access credantials contact with [**zeo.org**](https://zeo.org/contact-us/))
- Google Credantials. (have to enabled Drive API V3 on the account)
- Google Access Token. (to access to Google Drive as an user)
- Go v1.15

#### How to set up

- Copy `env.sample` to `.env`.  
- Update secret values in `.env`.

#### Usage at local

There is an easy way to test lambda projects at the local, `lambci`.

```shell
go build -o carbon && docker run --rm -v "$PWD":/var/task:ro,delegated lambci/lambda:go1.x carbon '{"HTTPMethod": "POST", "QueryStringParameters": {"type": "url", "format": "sheet", "country": "tr", "language": "tr"},"Body": "{\"values\": [{\"value\": \"https://boratanrikulu.dev\/contact\"}] }"}' && rm carbon
```

#### Run tests

To run all tests;
```shell
go test ./...
```

To run specific test;
```shell
go test ./services -run TestConvertURLResultToExcel -v 
```

#### How to deploy to AWS Lambda

Build.
```shell
go build -o carbon && zip deploy.zip carbon
```

Create the function.
```shell
aws lambda create-function --function-name CarbonLambda --handler carbon --runtime go1.x --role  arn:aws:iam::<account-id>:role/<role> --zip-file fileb://./deploy.zip --tracing-config Mode=Active
```

If you need to update the function, take a build and run this command.
```shell
aws lambda update-function-code --function-name CarbonLambda --zip-file fileb://./deploy.zip
```

## Credits

| [<img src="https://pbs.twimg.com/profile_images/935883931416657920/8HBYzvY7_400x400.jpg" width="100px;"/>](https://twitter.com/mertazizoglu) <br/> [T. Mert Azizoğlu](https://twitter.com/mertazizoglu)<br/><sub>Idea By</sub><br/> | [<img src="https://avatars3.githubusercontent.com/u/20258973?s=460&u=3147c97360ef8b5d64ef26c77077e1926a686356&v=4" width="100px;"/>](https://github.com/boratanrikulu) <br/>[Bora Tanrıkulu](https://github.com/boratanrikulu)<br/><sub>Developed By</sub><br/> |  
| - | - |

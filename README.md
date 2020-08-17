# ZEO.ORG - Carbon API

Carbon aims to find related URLs of given URLs or Keywords.  
It's mostly used to find alternative for 404 pages or SERP operations.  
It exports data in Excel or Google Sheets.

The API is served at AWS Lambda.

#### Requirements

- SERP API Credantials. (to access credantials contact with [**zeo.org**](https://zeo.org/contact-us/))
- Go v1.15

#### How to set up

- Copy `env.sample` to `.env`.  
- Update secret values in `.env`.
- Then you can run the app.  
  ```go run main.go```

#### Features

- Finds related results for the given input.  
- 2 input options; URLs and Keywords.  
	- URLs option; finds related 3 URLs.  
	  Mostly used to find options for 404 pages.  
	- Keywords option; finds related 10 URLs.  
	  The result includes title, desc and url for each input keywords.  
	  Mostly used for SERP.  
- Supports language specification for SERP.  
- 2 export options; Excel and Google Sheets.  
	- Suggested URLs are made bold for URLs options.  

#### Run tests

To run all tests;
```shell
go test ./...
```

To run specific test;
```shell
go test ./services -run TestConvertURLResultToExcel -v 
```

#### Endpoint

**URL:** /

**Request:**

- Type: **POST**
- Params: 
	- type `must`  
	  options: `keyword` or `url`.  
	- format `must`  
	  options: `excel` or `sheet`.
	- language `must`  
	  options: all languages supported by Google. 
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

- Type: 201
	- That means the data is created.
- Type: 400
	- That means the inputs are not valid.  
	  Check the error message.
- Type: 500
	- That means internal error occurs while creating the data.
- Type: 503
	- That means the service is not available.  
	  Try later.

Header and body;

- For excel;
	- Header  
		```
		Content-Disposition: attachment; filename="result.xlsx"
		Content-Length: 453646
		Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
		```
	- Body  
		`file`
- For sheet;  
	- Body  
		```
		 {
		     "sheetURL": "https://docs.google.com/spreadsheets/d/...",
		 }
		```

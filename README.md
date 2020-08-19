# ZEO.ORG - Carbon API

Carbon aims to find related URLs of given URLs or Keywords.  
It's mostly used to find alternative for 404 pages or SERP operations.  
It exports data in Excel or Google Sheets.

The API is served at AWS Lambda.

## Features

- Supports URL and Keyword options.
	- URLs option; finds related 3 URLs.  
	  Mostly used to find options for 404 pages.  
	- Keywords option; finds related 10 URLs.  
	  The result includes title, desc and url for each input keywords.  
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
	- type `must`  
	  options: `keyword` or `url`.  
	  note: `keyword` option is only available for internal users.
	- format `must`  
	  options: `excel` or `sheet`.
	- country `must`  
	  options: all countries supported by Google. 
	- langauge `must`  
	  options: all languages supported by Google.
	- accountName  
	- accountPassword  
- Header:
	- Accept  `must`  
	  If the format is `excel`  
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

- Type: 201
	- That means the data is created.
- Type: 400
	- That means the inputs are not valid.  
	  Check the error message.
- Type: 401
	- That means auth is not successful.
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

## Development

#### Requirements

- SERP API Credantials. (to access credantials contact with [**zeo.org**](https://zeo.org/contact-us/))
- Google Credantials. (have to enabled Drive API V3 on the account)
- Google Access Token. (to access to Google Drive as an user)
- Go v1.15

#### How to set up

- Copy `env.sample` to `.env`.  
- Update secret values in `.env`.
- Then you can run the app.  
  ```go run main.go```

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

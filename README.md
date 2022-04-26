# OSS-Score

We generate an activity and license score for open source projects based on their GitHub meta-data and the data of their dependencies. This data can be viewed in detail on our website, or inserted directly into the GitHub page via our Chrome extension

## Disclaimer
As of 4/17/22 GitHub has made all PATs contribute to the same over all rate limit of an account, which massivley reduced our throughput. We are temporarily not querying any dependencies and even regular queries are experiencing significant slowdowns. We are working on some alternative methods, but currently we are operating in at a significantly reduced capacity.

## Scoring

Scores are calculated based on github metadata and the metadata for the projects dependencies. Dependencies account for 25% of the score and the rest is the project itself. If we do not have the score for a dependency its score will be reported as max, with zero confidence.

### Confidence Rating

Sometimes we are unable to get certain metrics for a repository, or the dependency is yet to be calculated.
If this is the case our confidence rating for the provided score will decrease.

### Activity Score

Categories

* Issue Closure Rate

  Average time for an issue in the project to be closed. NOTE: This is only calculated based on the closed issues.
  * Weight: 25%
  * Linear Scale: 60 -- 0 closure in days

* Commit Cadence

  Average pace of commits in the project. Total number of commits divided by the query time frame.
  * Weight: 25%
  * Linear Scale: 0 -- 7 commits / week

* Contributors

  The number of unique users who have contributed to the project over the given query time frame.
  * Weight: 25%
  * Linear Scale: 0 -- 10 individual contributors

* Age of Last Release

  Time since the last release release.
  * Weight: 12.5%
  * Linear Scale: 52 -- 0 weeks since last release

* Release Cadence

  Average pace of releases in the project. Total number of releases divided by the query time frame.
  * Weight: 12.5%
  * Linear Scale: 0 -- 0.33 releases a month

* Pull Request Closure Rate (In progress)

  Average time for a pull request in the project to be closed. NOTE: This is only calculated based on the closed pull requests.

  * Weight: TBD
  * Linear Scale: TBD
### License Score

Direct mapping based on the license of the repo and the licenses of the dependencies.

Common Licenses:

* mit: 100
* gpl-3.0: 90
* unlicense: 100
* apache-2.0: 95
* bsd-3-clause: 85

The full specification can be found in [`licenseScoring.csv`](https://github.com/fonarevvichka/OSS-Score/blob/main/api/util/scores/licenseScoring.csv)

## Disclaimers / Limitations

* All metrics need to be pulled from GitHub, we try to cache as much as possible but sometimes queries can take a long time if we have no data cached. New repos can take anywhere from 30 seconds to a few minutes. Please be patient.

* We have an aritificial shelf life for our data. Data is considiered 'in-date' if is is less than three days old. As such some metrics may not quite line up with what you see on the repo homepage.

## Components

### API

The API is a completley serverless and uses AWS' API Gateway, SQS, and Lambda functions.

![OSS-Score API drawio](https://user-images.githubusercontent.com/14360853/163482102-a3ef41a7-a3f1-4da6-bcc7-35f610bf6cc3.png)

We use serverless for orchestration and deployment.
### DB

We use MongoDB since the large document sizes and NO-SQL structure suited us well. Hosted on Atlas beta serverless deployment.
### Website

React frontend hosted on Heroku.

### Chrome extension

The chrome extension inserts the scores directly into repo homepages in sidebar.
Once the extension is installed, it will retrieve to the score for that repo and insert the score.
If the score has not yet been calculated or is out of date, then a calculation request button will be shown instead.

NOTE: For now all chrome extension score calculations default to a 12 month time period.

#### Manual Installation Instructions

1. Navigate to the extensions page in your Chrome/Brave settings, usually `chrome://extensions`, and enable *developer mode*
2. Clone /  Download the repository.
3. Select load unpacked and then select the `OSS-Score/extension` directory
4. Enable the extension and watch the magic happen

### Deploy Yourself

We welcome you to deploy this project yourself!

#### Backend
To deploy the backend you will need to deploy two components, the API itself and the Database

#### DB
We host on Atlas but you are welcome to use whatever hosting/managing service you like. If using Atlas you can use the following instructions.
  1. Create an account with MongoDB and create your database. NOTE: Although MongoDB Atlas has a free tier we found that it often gets throttled to the point of unusability
  2. Once you have created your database make sure to allow access from any IP address and generate a `x509` certificate to be used by our lambda functions later.
  3. Name the certificate `mongo_cert.pem` and place it in `OSS-Score/api/util`

#### API
We use serverless as our deployment/orchestration tool and this will handle the majority of the deployment for you.
  1. Install and set up serverless as a global npm utility on your machine, make sure to add the correct AWS credentials.
  2. Navigate to `OSS-Score/api`. And run a build with `make`
  3. Run `sls deploy -s <env>`, the `-s` options lets you specify what environment you want to deploy to (the default is `dev`)

#### Enviroment Variables
DO NOT skip this section, if you do not set these up properly nothing will work.

You will need a single GitHub personal access token or PAT.
[Instructions](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token).
* The only permissions this PAT needs is `public_repo`. But if you want to be able to score your private repos then give them that additional permission.
* This PAT needs to be stored in your enviroment under `GIT_PAT`

There are a number of other enviroment variables that are all stored in `OSS-Score/api/vars.yaml`. The default values are all set but these allow you to customize the `shelf_life` of your data as well as the naming of different queues and database naming scheme.

### Frontend

#### Chrome Extension
Please see the above Chrome Extension section to see how to load the extension yourself.

To configure it to use your own version of the deployed API, edit the `basePath` constant at the top of `OSS-Score/extension/extension.js`

#### Website
The website is run as react app.
1. Navigate to `OSS-Score/website/oss-score-site` and run `npm install`
2. Start the server with `npm start`

To deploy to heroku follow their instructions and make sure to use `mars/create-react-app` buildpack


## Endpoints
All paths for an API Gateway will have the following prefix: `https://<id>.execute-api.us-east-2.amazonaws.com/<env></env>/`

See full Swagger definition in `OSS-Score/api/oss-score.yaml`

### getScore
* `OSS-Score/api/getScore/getScore.go`
* /catalog/{catalog}/owner/{owner}/name/{name}/type/{type}
* {catalog} is an enum
  * github
  * gitlab (Not yet supported)
* {type} is an enum
    * `license`
    * `activity`
* Query parameters
    * timeFrame: integer

### getMetric
* `OSS-Score/api/getMetric/getMetric.go`
* /catalog/{catalog}/owner/{owner}/name/{name}/metric/{metric}
* {catalog} is an enum
  * github
  * gitlab (Not yet supported)
* metric is an enum
  * all
  * stars
  * releaseCadence
  * ageLastRelease
  * commitCadence
  * contributors
  * issueClosureTime
  * repoActivityScore
  * dependencyActivityScore
  * repoLicenseScore
  * dependencyLicenseScore

* Query parameters
    * timeFrame: integer

### queryRepositoryScore
* `OSS-Score/api/queryRepository/handler/queryRepoHandler.go`
* /catalog/{catalog}/owner/{owner}/name/{name}
* {catalog} is an enum
  * github
  * gitlab (Not yet supported)
* Post BODY
  * timeFrame: integer string

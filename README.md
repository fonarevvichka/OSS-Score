# OSS-Score

## Scoring

Scores are calculated based on github metadata and the metadata for the projects dependencies. Dependencies account for 25% of the score and the rest is the project itself. If we do not have the score for a dependency its score will be reported as max, with zero confidence.

### Confidence Rating

Sometimes we are unable to get certain metrics for a repository, or the dependency is yet to be calculated.
If this is the case our confidence rating for the provided score will decrease.

### Activity Score

Categories

* Issue Closure Rate
  * Weight: 25%
  * Linear Scale: 176 -- 0 closure in days

* Commit Cadence
  * Weight: 25%
  * Linear Scale: 0 -- 2 commits / week
* Contributors
  * Weight: 25%
  * Linear Scale: 0 -- 10 individual contributors
* Age of Last Release
  * Weight: 12.5%
  * Linear Scale: 26 -- 0 weeks since last release
* Release Cadence
  * Weight: 12.5%
  * Linear Scale: 0 -- 0.33 releases a month

### License Score

Direct mapping based on the license of the repo and the licenses of the dependencies.

Common Licenses:

* mit: 100
* gpl-3.0: 90
* unlicense: 100
* apache-2.0: 95
* bsd-3-clause: 85

The full specification can be found in [`licenseScores.txt`](https://github.com/fonarevvichka/OSS-Score/blob/main/api/util/scores/licenseScores.txt)

## Components

### API

stuff stuff stuff

### Website

stuff stuff stuff

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
TODO

## Disclaimer / Limitations

TODO

## Endpoints

TODO

* getScore
* getMetric
* queryRepositoryScore

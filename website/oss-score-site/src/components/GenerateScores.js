import React from 'react'

import './GenerateScores.css';

export default function GenerateScores() {
    return (
        <div className='scores-container'>
            <header className='main-header'>How We Generate Scores</header>
            <header>Scores</header>
            <div className='sub-container'>
                &emsp; OSS-Score is driven to calculate scores and metrics that are useful to developers.
                 We provide developers with a comprehensive activity and license score, that is calculated from a wide variety of GitHub metadata for that repository and its dependencies. 
                 The project itself accounts for 75% of the score and the dependencies account for the other 25%. 
                 If we do not have the score for a dependency its score will be reported as max, with zero confidence. 
            </div>
            <header>Confidence Ratings</header>
            <div className='sub-container'>
                &emsp; Sometimes we are unable to get certain metrics for a repository, or the dependency is yet to be calculated.
                 If this is the case our confidence rating for the provided score will decrease.
                 Otherwise, we will always be 100% confident in our calculations.
            </div>
            <header>Calculating Activity Score</header>
            <div className='sub-container'>
                <div className="metric-desciption-container odd-metrics">
                    <div className='metric-title'>Issue Closure Rate</div>
                    <div className='metric-description'>
                        &emsp; Issue Closure Rate is the average time for an issue in the project to be closed. 
                            It is important to note that this metric is only calculated based on the closed issues 
                        <ul>
                            <li className='metric-weight'>Weight: 25%</li>
                            <li className='metric-linear-scale'>Linear Scale: 176 -- 0 closure in days</li>
                        </ul>
                    </div>
                </div>
                <div className="metric-desciption-container even-metrics">
                    <div className='metric-title'>Commit Cadence</div>
                    <div className='metric-description'>
                        &emsp; Commit Cadence is the average pace of commits in the project. 
                            More specifically, it is the total number of commits divided by the query time frame.
                        <ul>
                            <li className='metric-weight'>Weight: 25%</li>
                            <li className='metric-linear-scale'>Linear Scale: 0 -- 2 commits / week</li>
                        </ul>
                    </div>
                </div>
                <div className="metric-desciption-container odd-metrics">
                    <div className='metric-title'>Contributors</div>
                    <div className='metric-description'>
                        &emsp; The number of unique users who have contributed to the project over the given query time frame.
                        <ul>
                            <li className='metric-weight'>Weight: 25%</li>
                            <li className='metric-linear-scale'>Linear Scale: 0 -- 10 individual contributors</li>
                        </ul>
                    </div>
                </div>
                <div className="metric-desciption-container even-metrics">
                    <div className='metric-title'>Age of Last Release</div>
                    <div className='metric-description'>
                        &emsp; Age of Last Release is simply the time since the last release.
                        <ul>
                            <li className='metric-weight'>Weight: 12.5%</li>
                            <li className='metric-linear-scale'>Linear Scale: 26 -- 0 weeks since last release</li>
                        </ul>
                    </div>
                </div>
                <div className="metric-desciption-container odd-metrics">
                    <div className='metric-title'>Release Cadence</div>
                    <div className='metric-description'>
                        Release Cadence is the average pace of releases in the project.
                            More specifically, it is the total number of releases divided by the query time frame.
                        <ul>
                            <li className='metric-weight'>Weight: 12.5%</li>
                            <li className='metric-linear-scale'>Linear Scale: 0 -- 0.33 releases a month</li>
                        </ul>
                    </div>
                </div>
                <div className="metric-desciption-container even-metrics" style={{display:"none"}}>
                    <div className='metric-title'>Pull Request Closure Rate</div>
                    <div className='metric-description'>
                        &emsp; Pull Request Closure Rate is the average time for a pull request in the project to be closed.
                            It is important to note that this is only calculated based on the closed pull requests.
                        <ul>
                            <li className='metric-weight'>Weight: TBD</li>
                            <li className='metric-linear-scale'>Linear Scale: TBD</li>
                        </ul>
                    </div>
                </div>
            </div>
            <header>Calculating License Score</header>
            <div className='sub-container'>
                &emsp; License Score is a direct mapping based on the license of the repo and the licenses of the dependencies.
                <br></br><br></br>
                <div className="license-container odd-metrics">
                    <div className='license-title'>Common Licenses and Their Scores</div>
                    <div className='license-description'>
                        <ul>
                            <li className='license-example'>mit: 100</li>
                            <li className='license-example'>gpl-3.0: 90</li>
                            <li className='license-example'>unlicense: 100</li>
                            <li className='license-example'>apache-2.0: 95</li>
                            <li className='license-example'>bsd-3-clause: 85</li>
                        </ul>
                        <div className="license-file">The full specification can be found in 
                            <a href="https://github.com/fonarevvichka/OSS-Score/blob/main/api/util/scores/licenseScoring.csv" 
                            target="_blank" rel="noreferrer">licenseScoring.csv</a>
                        </div>
                    </div>
                </div>
            </div>


        </div>
  )
}

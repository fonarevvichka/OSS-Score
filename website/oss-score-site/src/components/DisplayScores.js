import React from 'react'

import './DisplayScores.css';


const DisplayScores = () => {
    return (
        <div class="repo1-stats">
            <div class="basic-info-display">
                <div class="basic-info-title">Name</div>
                <div class="basic-info" id="repoName1">None</div>
                <div class="basic-info-title">Author</div>
                <div class="basic-info" id="repoAuthor1">None</div>
            </div>
            {/* <div class="repo-header">
                <div class="metric-container-title">Overall Score</div>
                <div class="metric" id="overallScore1">0</div>
                <div class="confidence" id="overallConfScore1">Confidence: 0</div>
            </div> */}
            <div class="metrics-display">
                <div class="metric-category">Activity Scores</div>
                <div class="metric-container">
                    <div class="metric-container-title">Activity Score</div>
                    <div class="metric" id="repoActivityScore1">0</div>
                    <div class="confidence" id="repoActivityConfScore1">Confidence: 0</div>
                </div>
                <div class="metric-container">
                    <div class="metric-container-title">Issue Closure Time</div>
                    <div class="metric" id="issueClosureTime1">0</div>
                    <div class="confidence" id="issueClosureTimeConf1">Confidence: 0</div>
                </div>

                <div class="metric-container">
                    <div class="metric-container-title">Commit Cadence</div>
                    <div class="metric" id="commitCadence1">0</div>
                    <div class="confidence" id="commitCadenceConf1">Confidence: 0</div>
                </div>


                <div class="metric-container">
                    <div class="metric-container-title">Release Score</div>
                    {/* <div class="metric" id="releaseScore1">0</div>
                    <div class="confidence" id="releaseConfScore1">Confidence: 0</div> */}
                    <div class="submetric-container">
                        <div class="submetric-container-title">Age Last Release</div>
                        <div class="metric" id="ageLastRelease1">0</div>
                        <div class="confidence" id="ageLastReleaseConf1">Confidence: 0</div>
                    </div>
                    <div class="submetric-container">
                        <div class="submetric-container-title">Release Cadence</div>
                        <div class="metric" id="releaseCadence1">0</div>
                        <div class="confidence" id="releaseCadenceConf1">Confidence: 0</div>
                    </div>
                </div>

                <div class="metric-container">
                    <div class="metric-container-title">Contributors</div>
                    <div class="metric" id="contributors1">0</div>
                    <div class="confidence" id="contributorsConf1">Confidence: 0</div>
                </div>
            </div>
            <div class="repo-licence-score">
                <div class="metric-category">License Scores</div>
                <div class="metric-container">
                    <div class="metric-container-title">License Score</div>
                    <div class="metric" id="repoLicenseScore1">0</div>
                    <div class="confidence" id="repoLicenseConfScore1">Confidence: 0</div>
                </div>
            </div>
            <div class="repo-dependency-score">
                <div class="metric-category">Dependency Scores</div>
                <div class="metric-container">
                    <div class="metric-container-title">Dependency Activity Score</div>
                    <div class="metric" id="dependencyActivityScore1">0</div>
                    <div class="confidence" id="dependencyActivityConfScore1">Confidence: 0</div>
                </div>
                <div class="metric-container">
                    <div class="metric-container-title">Dependency License Score</div>
                    <div class="metric" id="dependencyLicenseScore1">0</div>
                    <div class="confidence" id="dependencyLicenseConfScore1">Confidence: 0</div>
                </div>
            </div>
            <div class="repo-stars">
                <div class="metric-category">Stars</div>
                <div class="metric-container">
                    <div class="metric-container-title">Stars</div>
                    <div class="metric" id="stars1">0</div>
                    <div class="confidence" id="stars1">Confidence: 0</div>
                </div>
            </div>
        </div>
    );
}

export default DisplayScores;
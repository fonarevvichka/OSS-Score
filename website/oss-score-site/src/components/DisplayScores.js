import React from 'react'

import './DisplayScores.css';
import './Homepage.css';


const DisplayScores = (owner, name, metrics) => {
    // round scores
    
    for (var key in metrics) {
        if (key != "message") {
            metrics[key].metric = Math.round(metrics[key].metric * 100) / 100
        }
    }

    return (
        '<div class="repo-stats">\
            <div class="basic-info-display">\
                <div class="basic-info-title">Name</div>\
                <div class="basic-info" id="repoName">' + name + '</div>\
                <div class="basic-info-title">Author</div>\
                <div class="basic-info" id="repoAuthor">' + owner + '</div>\
            </div>\
            <div class="metrics-display">\
                <div class="metric-category">Activity Scores</div>\
                <div class="metric-container">\
                    <div class="metric-container-title">Activity Score</div>\
                    <div class="metric" id="repoActivityScore">' + metrics.repoActivityScore.metric + '</div>\
                    <div class="confidence" id="repoActivityConfScore">Confidence: ' + metrics.repoActivityScore.confidence + '</div>\
                </div>\
                <div class="metric-container">\
                    <div class="metric-container-title">Issue Closure Time</div>\
                    <div class="metric" id="issueClosureTime">' + metrics.issueClosureTime.metric + '</div>\
                    <div class="confidence" id="issueClosureTimeConf">Confidence: ' + metrics.issueClosureTime.confidence + '</div>\
                </div>\
                <div class="metric-container">\
                    <div class="metric-container-title">Commit Cadence</div>\
                    <div class="metric" id="commitCadence">' + metrics.commitCadence.metric + '</div>\
                    <div class="confidence" id="commitCadenceConf">Confidence: ' + metrics.commitCadence.confidence + '</div>\
                </div>\
                <div class="metric-container">\
                    <div class="metric-container-title">Release Score</div>\
                    <div class="submetric-container">\
                        <div class="submetric-container-title">Age Last Release</div>\
                        <div class="metric" id="ageLastRelease">' + metrics.ageLastRelease.metric + '</div>\
                        <div class="confidence" id="ageLastReleaseConf">Confidence: ' + metrics.ageLastRelease.confidence + '</div>\
                    </div>\
                    <div class="submetric-container">\
                        <div class="submetric-container-title">Release Cadence</div>\
                        <div class="metric" id="releaseCadence">' + metrics.releaseCadence.metric + '</div>\
                        <div class="confidence" id="releaseCadenceConf">Confidence: ' + metrics.releaseCadence.confidence + '</div>\
                    </div>\
                </div>\
                <div class="metric-container">\
                    <div class="metric-container-title">Contributors</div>\
                    <div class="metric" id="contributors">' + metrics.contributors.metric + '</div>\
                    <div class="confidence" id="contributorsConf">Confidence: ' + metrics.contributors.confidence + '</div>\
                </div>\
            </div>\
            <div class="repo-licence-score">\
                <div class="metric-category">License Scores</div>\
                <div class="metric-container">\
                    <div class="metric-container-title">License Score</div>\
                    <div class="metric" id="repoLicenseScore">' + metrics.repoLicenseScore.metric + '</div>\
                    <div class="confidence" id="repoLicenseConfScore">Confidence: ' + metrics.repoLicenseScore.confidence + '</div>\
                </div>\
            </div>\
            <div class="repo-dependency-score">\
                <div class="metric-category">Dependency Scores</div>\
                <div class="metric-container">\
                    <div class="metric-container-title">Dependency Activity Score</div>\
                    <div class="metric" id="dependencyActivityScore">' + metrics.dependencyActivityScore.metric + '</div>\
                    <div class="confidence" id="dependencyActivityConfScore">Confidence: ' + metrics.dependencyActivityScore.confidence + '</div>\
                </div>\
                <div class="metric-container">\
                    <div class="metric-container-title">Dependency License Score</div>\
                    <div class="metric" id="dependencyLicenseScore">' + metrics.dependencyLicenseScore.metric + '</div>\
                    <div class="confidence" id="dependencyLicenseConfScore">Confidence: ' + metrics.dependencyLicenseScore.confidence + '</div>\
                </div>\
            </div>\
            <div class="repo-stars">\
                <div class="metric-category">Stars</div>\
                <div class="metric-container">\
                    <div class="metric-container-title">Stars</div>\
                    <div class="metric" id="stars">' + metrics.stars.metric + '</div>\
                    <div class="confidence" id="stars">Confidence: ' + metrics.stars.confidence + '</div>\
                </div>\
            </div>\
        </div>'
    );
}

export default DisplayScores;
import './DisplayScores.css';
import './Homepage.css';


const DisplayScores = (owner, name, metrics) => {
    // round scores
    
    for (var key in metrics) {
        if (key !== "message") {
            metrics[key].metric = Math.round(metrics[key].metric * 100) / 100
        }
    }

    return (
        '<div class="repo-stats"> \n' +
        '<div class="basic-info-display"> \n' +
                '<div class="basic-info-title">Name</div>\n' +
                '<div class="basic-info" id="repoName">' + name + '</div>\n' +
                '<div class="basic-info-title">Author</div>\n' +
                '<div class="basic-info" id="repoAuthor">' + owner + '</div>\n' +
            '</div>\n' +
            '<div class="metrics-display">\n' +
                '<div class="metric-category">Activity Scores</div>\n' +
                '<div class="metric-container">\n' +
                    '<div class="metric-container-title">Activity Score</div>\n' +
                    '<div class="metric" id="repoActivityScore">' + metrics.repoActivityScore.metric + '</div>\n' +
                    '<div class="confidence" id="repoActivityConfScore">Confidence: ' + metrics.repoActivityScore.confidence + '</div>\n' +
                '</div>\n' +
                '<div class="metric-container">\n' +
                    '<div class="metric-container-title">Issue Closure Time</div>\n' +
                    '<div class="metric" id="issueClosureTime">' + metrics.issueClosureTime.metric + '</div>\n' +
                    '<div class="confidence" id="issueClosureTimeConf">Confidence: ' + metrics.issueClosureTime.confidence + '</div>\n' +
                '</div>\n' +
                '<div class="metric-container">\n' +
                    '<div class="metric-container-title">Commit Cadence</div>\n' +
                    '<div class="metric" id="commitCadence">' + metrics.commitCadence.metric + '</div>\n' +
                    '<div class="confidence" id="commitCadenceConf">Confidence: ' + metrics.commitCadence.confidence + '</div>\n' +
                '</div>\n' +
                '<div class="metric-container">\n' +
                    '<div class="metric-container-title">Release Score</div>\n' +
                    '<div class="submetric-container">\n' +
                        '<div class="submetric-container-title">Age Last Release</div>\n' +
                        '<div class="metric" id="ageLastRelease">' + metrics.ageLastRelease.metric + '</div>\n' +
                        '<div class="confidence" id="ageLastReleaseConf">Confidence: ' + metrics.ageLastRelease.confidence + '</div>\n' +
                    '</div>\n' +
                    '<div class="submetric-container">\n' +
                        '<div class="submetric-container-title">Release Cadence</div>\n' +
                        '<div class="metric" id="releaseCadence">' + metrics.releaseCadence.metric + '</div>\n' +
                        '<div class="confidence" id="releaseCadenceConf">Confidence: ' + metrics.releaseCadence.confidence + '</div>\n' +
                    '</div>\n' +
                '</div>\n' +
                '<div class="metric-container">\n' +
                    '<div class="metric-container-title">Contributors</div>\n' +
                    '<div class="metric" id="contributors">' + metrics.contributors.metric + '</div>\n' +
                    '<div class="confidence" id="contributorsConf">Confidence: ' + metrics.contributors.confidence + '</div>\n' +
                '</div>\n' +
            '</div>\n' +
            '<div class="repo-licence-score">\n' +
                '<div class="metric-category">License Scores</div>\n' +
                '<div class="metric-container">\n' +
                    '<div class="metric-container-title">License Score</div>\n' +
                    '<div class="metric" id="repoLicenseScore">' + metrics.repoLicenseScore.metric + '</div>\n' +
                    '<div class="confidence" id="repoLicenseConfScore">Confidence: ' + metrics.repoLicenseScore.confidence + '</div>\n' +
                '</div>\n' +
            '</div>\n' +
            '<div class="repo-dependency-score">\n' +
                '<div class="metric-category">Dependency Scores</div>\n' +
                '<div class="metric-container">\n' +
                    '<div class="metric-container-title">Dependency Activity Score</div>\n' +
                    '<div class="metric" id="dependencyActivityScore">' + metrics.dependencyActivityScore.metric + '</div>\n' +
                    '<div class="confidence" id="dependencyActivityConfScore">Confidence: ' + metrics.dependencyActivityScore.confidence + '</div>\n' +
                '</div>\n' +
                '<div class="metric-container">\n' +
                    '<div class="metric-container-title">Dependency License Score</div>\n' +
                    '<div class="metric" id="dependencyLicenseScore">' + metrics.dependencyLicenseScore.metric + '</div>\n' +
                    '<div class="confidence" id="dependencyLicenseConfScore">Confidence: ' + metrics.dependencyLicenseScore.confidence + '</div>\n' +
                '</div>\n' +
            '</div>\n' +
            '<div class="repo-stars">\n' +
                '<div class="metric-category">Stars</div>\n' +
                '<div class="metric-container">\n' +
                    '<div class="metric-container-title">Stars</div>\n' +
                    '<div class="metric" id="stars">' + metrics.stars.metric + '</div>\n' +
                    '<div class="confidence" id="stars">Confidence: ' + metrics.stars.confidence + '</div>\n' +
                '</div>\n' +
            '</div>\n' +
        '</div>'
    );
}

export default DisplayScores;
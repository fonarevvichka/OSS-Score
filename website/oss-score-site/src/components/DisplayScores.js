import './DisplayScores.css';
import './Homepage.css';


const getMetricDisplay = (metricScore, barDisplay) => {
    // round metric 
    metricScore = Math.round(metricScore * 100) / 100
    if (barDisplay) {
        // TODO: create bar graphic
        return toString(metricScore)
    } else {
        // just display raw score
        return toString(metricScore)
    }
}

const getMetricContainer = (metricName, metricScore, metricConf, barDisplay, highlight) => {

    if (highlight) {
        // highlight score in green
        return '<div class="metric-container" style="border-color: green;">\n' +
            '<div class="metric-container-title">' + metricName + '</div>\n' +
            '<div class="metric">' + getMetricDisplay(metricScore, barDisplay) + '</div>\n' +
            '<div class="confidence">Confidence: ' + metricConf + '</div>\n' +
            '</div>\n'
    } else {
        // do not highlight score
        return '<div class="metric-container">\n' +
            '<div class="metric-container-title">' + metricName + '</div>\n' +
            '<div class="metric">' + getMetricDisplay(metricScore, barDisplay) + '</div>\n' +
            '<div class="confidence">Confidence: ' + metricConf + '</div>\n' +
            '</div>\n'

    }
}

// subMetrics is list of tuples (nameOfMetric, MetricScore, MetricConf, barDisplay, highlight)
const getMetricContainerWSubContainers = (metricName, subMetrics, barDisplay, highlight) => {
    let subcontainers = ""

    for (let i = 0; i < subMetrics.length; i++) {
        if (subMetrics[i][4]) {
            subcontainers += '<div class="submetric-container" style="border-color: green;>\n' +
                '<div class="submetric-container-title">' + subMetrics[i][0] + '</div>\n' +
                '<div class="metric" id="ageLastRelease">' + getMetricDisplay(subMetrics[i][1], subMetrics[i][3]) + '</div>\n' +
                '<div class="confidence" id="ageLastReleaseConf">Confidence: ' + subMetrics[i][2] + '</div>\n' +
                '</div>\n'
        } else {
            subcontainers += '<div class="submetric-container">\n' +
                '<div class="submetric-container-title">' + subMetrics[i][0] + '</div>\n' +
                '<div class="metric" id="ageLastRelease">' + getMetricDisplay(subMetrics[i][1], subMetrics[i][3]) + '</div>\n' +
                '<div class="confidence" id="ageLastReleaseConf">Confidence: ' + subMetrics[i][2] + '</div>\n' +
                '</div>\n'
        }
    }


    return '<div class="metric-container">\n' +
        '<div class="metric-container-title">' + metricName + '</div>' + subcontainers + '</div>'
}

// activityScore, licenseScore, stars, contributors are tuples (metricScore, confidence, highlight)
const getBasicInfoDisplay = (owner, name, activityScore, licenseScore, stars, contributors) => {
    let result = '<div class="basic-info-display"> \n'
    '<div class="basic-info" id="repoOwnerName">' + owner + '/' + name + '</div>'

    if (activityScore[2]) {
        result += '<div class="basic-info" id="activityScore" style="border-color: green;> Activity Score: ' + activityScore[0] + '/10</div>'
    } else {
        result += '<div class="basic-info" id="activityScore"> Activity Score: ' + activityScore[0] + '/10</div>'
    }
    result += '<div class="basic-info-conf" id="activityScoreConf"> Confidence: ' + activityScore[1] + '</div>'

    if (licenseScore[2]) {
        result += '<div class="basic-info" id="licenseScore" style="border-color: green;> License Score: ' + licenseScore + '/10</div>'
    } else {
        result += '<div class="basic-info" id="licenseScore"> License Score: ' + licenseScore[0] + '/10</div>'
    }
    result += '<div class="basic-info-conf" id="licenseScoreConf"> Confidence: ' + licenseScore[1] + '</div>'

    if (stars[2]) {
        result += '<div class="basic-info" id="stars" style="border-color: green;>Stars: ' + stars[0] + '</div>'
    } else {
        result += '<div class="basic-info" id="stars">Stars: ' + stars[0] + '</div>'
    }
    result += '<div class="basic-info-conf" id="stars">Confidence: ' + stars[1] + '</div>'

    if (contributors[2]) {
        result += '<div class="basic-info" id="contributors" style="border-color: green;>Contributors: ' + contributors[0] + '</div>'
    } else {
        result += '<div class="basic-info" id="contributors">Contributors: ' + contributors[0] + '</div>'
    }
    result += '<div class="basic-info-conf" id="contributors">Confidence: ' + contributors[1] + '</div>'



    // close div
    result += '</div >'
    return result
}

// metrics is array of json objects
const DisplayScores1 = (owner, name, metrics) => {
    let result = '<div class="repo-stats">'

    result += getBasicInfoDisplay()

    result += '<div class="metrics-display">'
    result += '<div class="metric-category">Activity Scores</div>'
    result += getMetricContainer('Issue Closure Time')
    result += getMetricContainer('Commit Cadence')
    result += getMetricContainerWSubContainers('Release Score')
    result += '</div >'
    result += '</div >'

    result += '<div class="repo-dependency-score">'
    result += '<div class="metric-category">Dependency Scores</div>'
    result += getMetricContainer('Dependency Activity Score')
    result += getMetricContainer('Dependency License Score')

    result += '</div >'
    result += '</div >'

    result += '</div >'
}


const DisplayScores = (owner, name, metrics) => {
    // round scores
    
    // rounding scores will be done in getMetricDisplay
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
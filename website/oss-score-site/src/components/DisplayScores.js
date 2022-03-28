import './DisplayScores.css';
import './Homepage.css';


const getMetricDisplay = (metricScore, barDisplay) => {
    // round metric 
    metricScore = Math.round(metricScore * 100) / 100
    if (barDisplay) {
        // TODO: create bar graphic
        return metricScore
    } else {
        // just display raw score
        return metricScore
    }
}

const getMetricContainer = (metricName, metric, barDisplay) => {

    if (metric.highlight) {
        // highlight score in green
        return '<div class="metric-container" style="color: green;">\n' +
            '<div class="metric-container-title">' + metricName + '</div>\n' +
            '<div class="metric">' + getMetricDisplay(metric.metric, barDisplay) + '</div>\n' +
            '<div class="confidence">Confidence: ' + metric.confidence + '</div>\n' +
            '</div>\n'
    } else {
        // do not highlight score
        return '<div class="metric-container">\n' +
            '<div class="metric-container-title">' + metricName + '</div>\n' +
            '<div class="metric">' + getMetricDisplay(metric.metric, barDisplay) + '</div>\n' +
            '<div class="confidence">Confidence: ' + metric.confidence + '</div>\n' +
            '</div>\n'

    }
}

// subMetrics is list of tuples (nameOfMetric, MetricScore, MetricConf, barDisplay, highlight)
const getMetricContainerWSubContainers = (metricName, subMetrics) => {
    let subcontainers = ""

    //    (nameOfMetric, Metric, barDisplay)

    for (let i = 0; i < subMetrics.length; i++) {
        if (subMetrics[i][1].highlight) {
            subcontainers += '<div class="submetric-container" style="color: green;">\n' +
                '<div class="submetric-container-title">' + subMetrics[i][0] + '</div>\n' +
                '<div class="metric">' + getMetricDisplay(subMetrics[i][1].metric, subMetrics[i][2]) + '</div>\n' +
                '<div class="confidence">Confidence: ' + subMetrics[i][1].confidence + '</div>\n' +
                '</div>\n'
        } else {
            subcontainers += '<div class="submetric-container">\n' +
                '<div class="submetric-container-title">' + subMetrics[i][0] + '</div>\n' +
                '<div class="metric">' + getMetricDisplay(subMetrics[i][1].metric, subMetrics[i][2]) + '</div>\n' +
                '<div class="confidence">Confidence: ' + subMetrics[i][1].confidence + '</div>\n' +
                '</div>\n'
        }
    }


    return '<div class="metric-container">\n' +
        '<div class="metric-container-title">' + metricName + '</div>' + subcontainers + '</div>'
}

// activityScore, licenseScore, stars, contributors are tuples (metricScore, confidence, highlight)
const getBasicInfoDisplay = (owner, name, activityScore, licenseScore, stars, contributors) => {
    let result = '<div class="basic-info-display"> \n' +
    '<div class="basic-info" id="repoOwnerName">' + owner + '/' + name + '</div>'

    if (activityScore.highlight) {
        result += '<div class="basic-info" id="activityScore" style="color: green";> Activity Score: ' + getMetricDisplay(activityScore.metric, false) + '/10</div>'
    } else {
        //result += '<div class="basic-info" id="activityScore" style="border-color: green;> Activity Score: ' + getMetricDisplay(activityScore.metric, false) + '/10</div>'
        result += '<div class="basic-info" id="activityScore"> Activity Score: ' + getMetricDisplay(activityScore.metric, false) + '/10</div>'
    }
    result += '<div class="basic-info-conf" id="activityScoreConf"> Confidence: ' + activityScore.confidence + '</div>'

    if (licenseScore.highlight) {
        result += '<div class="basic-info" id="licenseScore" style="color: green";> License Score: ' + getMetricDisplay(licenseScore.metric, false) + '/10</div>'
    } else {
        result += '<div class="basic-info" id="licenseScore"> License Score: ' + getMetricDisplay(licenseScore.metric, false) + '/10</div>'
    }
    result += '<div class="basic-info-conf" id="licenseScoreConf"> Confidence: ' + licenseScore.confidence + '</div>'

    if (stars.highlight) {
        result += '<div class="basic-info" id="stars" style="color: green";>Stars: ' + stars.metric + '</div>'
    } else {
        result += '<div class="basic-info" id="stars">Stars: ' + stars.metric + '</div>'
    }
    result += '<div class="basic-info-conf" id="stars">Confidence: ' + stars.confidence + '</div>'

    if (contributors.highlight) {
        result += '<div class="basic-info" id="contributors" style="color: green";>Contributors: ' + contributors.metric + '</div>'
    } else {
        result += '<div class="basic-info" id="contributors">Contributors: ' + contributors.metric + '</div>'
    }
    result += '<div class="basic-info-conf" id="contributors">Confidence: ' + contributors.confidence + '</div>'



    // close div
    result += '</div >'
    return result
}

const AddHighlightJSON = (metricsAll) => {
    // console.log(metricsAll)
    // Add logic so that we add highlight bool to each thing 
    if (metricsAll.length < 1) {
        alert("Scores Unavailable")
    }
    for (var key in metricsAll[0]) {
        //console.log(key)
        if (key !== 'message') {
            let maxOfMetric = 0;
            let metricArray = [];
            // Find max of metrics and store values for highlighting
            for (let i = 0; i < metricsAll.length; i++) {
                metricArray.push(metricsAll[i][key].metric)
                if (metricsAll[i][key].metric > maxOfMetric) {
                    maxOfMetric = metricsAll[i][key].metric;
                }
            }

            // Add highligt property metrics
            for (let i = 0; i < metricArray.length; i++) {
                if (metricArray[i] === maxOfMetric) {
                    // Highlight metric
                    metricsAll[i][key].highlight = true;
                } else {
                    // Don't highlight
                    metricsAll[i][key].highlight = false;
                }
            }
        }
    }

    return metricsAll
}


const DisplayScores = (metrics) => {
    console.log(metrics)

    //let {owner1, name1, metricsAll} = allmetrics

    let metricsAll = [] 
    for (let i = 0; i < metrics.length; i++) {
        metricsAll.push(metrics[i][2])
    }
    // console.log(metricsAll)
    metricsAll = AddHighlightJSON(metricsAll)
    // console.log(metricsAll)
    
    let result = ''

    for (let i = 0; i < metrics.length; i++) {
        result += '<div class="repo-stats">'

        // owner, name, activityScore, licenseScore, stars, contributors
        result += getBasicInfoDisplay(metrics[i][0], metrics[i][1], metricsAll[i].repoActivityScore,
            metricsAll[i].repoLicenseScore, metricsAll[i].stars, metricsAll[i].contributors)



            // metricName, metricScore, metricConf, barDisplay, highlight

        result += '<div class="metrics-display">'
        result += '<div class="metric-category">Activity Scores</div>'
        result += getMetricContainer('Issue Closure Time', metricsAll[i].issueClosureTime, true)
        result += getMetricContainer('Commit Cadence', metricsAll[i].commitCadence, true)

        let releaseMetrics = [['Release Cadence', metricsAll[i].releaseCadence, true], ['Age of Last Release', metricsAll[i].ageLastRelease, true]]

        result += getMetricContainerWSubContainers('Release Score', releaseMetrics)
        result += '</div >'
       // result += '</div >'

        
            // subMetrics is list of tuples (nameOfMetric, Metric, barDisplay)

        result += '<div class="repo-dependency-score">'
        result += '<div class="metric-category">Dependency Scores</div>'
        result += getMetricContainer('Dependency Activity Score', metricsAll[i].dependencyActivityScore, true)
        result += getMetricContainer('Dependency License Score', metricsAll[i].dependencyLicenseScore, true)

        result += '</div >'
        result += '</div >'

        result += '</div >'
    }

    return result
    
    let owner = ""
    let name = ""
    
    
    
    
    
    // round scores
    
    //console.log(metrics)

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


// metrics is array of json objects
const DisplayScores1 = (metrics) => {
    // Add logic so that we add highlight bool to each thing 
    alert("inside displayscores 1")
    console.log("inside display scores 1")
    if (metrics.length < 1) {
        alert("Scores Unavailable")
    }

    // Add attribute to json's to indicate which fields to highlight
    for (var key in metrics[0]) {
        if (key !== 'message') {
            let maxOfMetric = 0;
            let metricArray = [];
            // Find max of metrics and store values for highlighting
            for (let i = 0; i < metrics.length; i++) {
                metricArray.push(metrics[i].key.metric)
                if (metrics[i].key.metric > maxOfMetric) {
                    maxOfMetric = metrics[i].key.metric;
                }
            }

            // Add highligt property metrics
            for (let i = 0; i < metricArray.length; i++) {
                if (metricArray[i] === maxOfMetric) {
                    // Highlight metric
                    metrics[i].key.highlight = true;
                } else {
                    // Don't highlight
                    metrics[i].key.highlight = false;
                }
            }
        }
    }

    console.log(metrics)

    // let result = ''

    // for (let i = 0; i < metrics.length; i++) {
    //     result += '<div class="repo-stats">'

    //     // owner, name, activityScore, licenseScore, stars, contributors
    //     result += getBasicInfoDisplay(owner, name, )

    //     result += '<div class="metrics-display">'
    //     result += '<div class="metric-category">Activity Scores</div>'
    //     result += getMetricContainer('Issue Closure Time')
    //     result += getMetricContainer('Commit Cadence')
    //     result += getMetricContainerWSubContainers('Release Score')
    //     result += '</div >'
    //     result += '</div >'

    //     result += '<div class="repo-dependency-score">'
    //     result += '<div class="metric-category">Dependency Scores</div>'
    //     result += getMetricContainer('Dependency Activity Score')
    //     result += getMetricContainer('Dependency License Score')

    //     result += '</div >'
    //     result += '</div >'

    //     result += '</div >'
    // }

    // return result
}

export default DisplayScores;
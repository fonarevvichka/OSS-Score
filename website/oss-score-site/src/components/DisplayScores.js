import './DisplayScores.css';
import './Homepage.css';



// const showFile = (e) => {
//     e.preventDefault();
//     const reader = new FileReader();
//     reader.onload = (e) => {
//         const text = e.target.result;
//         console.log(text);
//     };
//     reader.readAsText(e.target.files[0]);
// };




// TODO: Add confidence to metric display 
const getMetricDisplay = (metricScore, barDisplay, metricMin, metricMax) => {
    // round metric 
    //metricScore.metric = Math.round(metricScore.metric * 100) / 100
    metricScore = Math.round(metricScore * 100) / 100
    if (barDisplay) {


        // TODO:
        //      make txt file and put metric mins and maxes, then load it in hash tables
        //      figure out where to load the hash tabl, app.js??
        //      finish styling for metric bar
        //          Add low and high txt
        //          Make black bar extend out of colored bar a bit
        //          put confidence under bar
        //          put score and name of score next to bar


        //return metricScore
        //return '<div class="bar"></div>'
        // return '<table> \n'+
        //     '<tr> \n' +
        //     '<td></td> \n' +
        //     '<td rowspan=2><div class="bar"></div></td> \n' +

        //     '</tr>\n' +
        //     '<tr>\n' +
        //     '<td rowspan=2><div class="pointer"></div></td>\n' +
        //     '</tr>\n' +
        //     '</table>'

        // style="right:' + metricScore'%

        // TODO: make txt file and put metric mins and maxes, then load into
        //       hash tables

        let metricPercentage = metricScore.metric / (metricMax - metricMin) 
        
        // Cap percentage at 100%
        if (metricPercentage > 100) {
            metricPercentage = 100
        }

        return '<div class="metric-num">' + metricScore.metric +
               '</div><div class="bar-display"> \n' +
                '<div class="low">LOW</div>\n' +
                '<div class="bar"></div>\n' +
                '<div class="high">HIGH</div>\n' +
                '<div class="pointer" style="left: ' + metricPercentage + '%;"></div>\n' +
               '</div>'

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
// Given an array of JSON objects, add highlight bool to highlight fields with max values
const AddHighlightJSON = (metricsAll) => {
    // Add highlight bool to each thing field of the JSON
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
}


const DisplayScores = (metrics) => {
    // Extract JSON objects from metrics
    let metricsAll = [] 
    for (let i = 0; i < metrics.length; i++) {
        if (metrics[i] != null) {
            metricsAll.push(metrics[i][2])
        } else {
            metricsAll.push(null)
        }
    }

    // Add fields to JSON for highlighting
    AddHighlightJSON(metricsAll)
    
    let result = ''

    // Create head to head display for each owner/name/metrics
    for (let i = 0; i < metrics.length; i++) {
        result += '<div class="repo-stats">'
        
        // if (metrics[i] == null) {
        //     result += '<div class="no-metrics">No metrics to display</div>'
        //     continue
        // }

        // owner, name, activityScore, licenseScore, stars, contributors
        result += getBasicInfoDisplay(metrics[i][0], metrics[i][1], metricsAll[i].repoActivityScore,
            metricsAll[i].repoLicenseScore, metricsAll[i].stars, metricsAll[i].contributors)

        result += '<div class="metrics-display">'
        result += '<div class="metric-category">Activity Scores</div>'
        result += getMetricContainer('Issue Closure Time', metricsAll[i].issueClosureTime, true)
        result += getMetricContainer('Commit Cadence', metricsAll[i].commitCadence, true)

        let releaseMetrics = [['Release Cadence', metricsAll[i].releaseCadence, true], ['Age of Last Release', metricsAll[i].ageLastRelease, true]]

        result += getMetricContainerWSubContainers('Release Score', releaseMetrics)
        result += '</div >'
        

        result += '<div class="repo-dependency-score">'
        result += '<div class="metric-category">Dependency Scores</div>'
        result += getMetricContainer('Dependency Activity Score', metricsAll[i].dependencyActivityScore, true)
        result += getMetricContainer('Dependency License Score', metricsAll[i].dependencyLicenseScore, true)

        result += '</div >'
        result += '</div >'

        result += '</div >'
    }

    return result
}

export default DisplayScores;
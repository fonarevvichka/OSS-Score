import './DisplayScores.css';
import './Homepage.css';

// note: For stats where lower numbers are better, min and max are switched
let MetricStats = {

    "Activity Score-tooltip": "Overall activity score based on github metadata",

    "License Score-tooltip": "Direct mapping based on the license of the repo",

    "Stars-tooltip": "The number of stars gazers for the repository",

    "Contributors-tooltip": "The number of unique users who have contributed to the project over the given query time frame",

    "Dependency Activity Score-tooltip": "Overall activity score based on github metadata for project dependencies",

    "Dependency License Score-tooltip": "Direct mapping based on the licenses of the dependencies",

    "Issue Closure Time-min": 176,
    "Issue Closure Time-max": 0,
    "Issue Closure Time-units": "days",
    "Issue Closure Time-tooltip": "Average time for an issue in the project to be closed. NOTE: This is only calculated based on the closed issues",

    "Commit Cadence-min": 0,
    "Commit Cadence-max": 2,
    "Commit Cadence-units": "commits/week",
    "Commit Cadence-tooltip": "Average pace of commits in the project. Total number of commits divided by the query time frame",

    "Release Cadence-min": 0,
    "Release Cadence-max": 0.33,
    "Release Cadence-units": "releases/month",
    "Release Cadence-tooltip": "Average pace of releases in the project. Total number of releases divided by the query time frame",

    "Age of Last Release-min": 26,
    "Age of Last Release-max": 0,
    "Age of Last Release-units": "weeks",
    "Age of Last Release-tooltip": "Time since the last release release",
}


// TODO: Highlighting better metric, when lower is bette
const getMetricDisplay = (metricScore, metricName, barDisplay, outOfTen) => {
    // round metric 
    metricScore.metric = Math.round(metricScore.metric * 100) / 100
    metricScore.confidence = Math.round(metricScore.confidence)

    let result = ''

    if (barDisplay) {
        result += '<div class="metric-container bar-container" '
    } else {
        result += '<div class="metric-container lg-container" '
    }



    if (metricScore.highlight) {
        result += 'style="background-color: #b0c4de;">'
    } else {
        result += 'style="background-color: #d3d3d3;">'
    }

    result += '<div class="tool-tip">i\n' +
        '<span class="tooltiptext">' + MetricStats[metricName + "-tooltip"] + '</span>\n' +
        '</div>\n' +
        '<div class="metric-container-title">' + metricName + '</div>\n' +
        '<div class="metrics">' 


    
    // if (metricScore.highlight) {
    //     result += '<div class="metric-container" style="background-color: #b0c4de;">\n' +
    //     '<div class="tool-tip">i\n' + 
    //     '<span class="tooltiptext">' + MetricStats[metricName+"-tooltip"] + '</span>\n' +
    //     '</div>\n' +
    //     '<div class="metric-container-title">' + metricName + '</div>\n' +
    //     '<div class="metrics">' 

    // } else {
    //     result += '<div class="metric-container" style="background-color: #d3d3d3;">\n' +
    //         '<div class="tool-tip">i\n' +
    //         '<span class="tooltiptext">' + MetricStats[metricName+"-tooltip"] + '</span>\n' +
    //         '</div>\n' +
    //         '<div class="metric-container-title">' + metricName + '</div>\n' +
    //         '<div class="metrics">'
    // }

    if (barDisplay) {

        let metricMin = MetricStats[metricName + "-min"];
        let metricMax = MetricStats[metricName + "-max"];

        let metricPercentage = ((metricScore.metric - metricMin) / (metricMax - metricMin)) * 100
        
        // Cap percentage at 0% and 100%
        if (metricPercentage > 100) {
            metricPercentage = 100
        } else if (metricPercentage < 0) {
            metricPercentage = 0
        }

        result += '<div class="metric-num">' + metricScore.metric + ' ' + MetricStats[metricName + "-units"] +
                  '</div><div class="bar-and-conf">\n' +  
                            '<div class="bar-display"> \n' +
                                '<div class="bar">\n' + 
                                    '<div class="low">Low</div>\n' + 
                                    '<div class="high">High</div>\n' +
                                '</div>\n' +
                                '<div class="pointer" style="left: ' + metricPercentage + '%;"></div>\n' +
                            '</div>\n' +
            '<div class="metric-confidence"> Confidence: ' + metricScore.confidence + '%</div></div>'

    } else {
        // display raw score and confidence
        result += '<div class="metric-num">' + metricScore.metric
        if (outOfTen) {
            result += '/10'
        }

        result += '</div><div class="metric-confidence">Confidence: ' + metricScore.confidence + '%</div>'

        // if (outOfTen) {
        //     // display metric out of 10
        //     result += '<div class="metric-num">' + metricScore.metric + '/10</div> \n' +
        //           '<div class="metric-confidence">Confidence: ' + metricScore.confidence + '%</div>'
        // } else {
        //     result += '<div class="metric-num">' + metricScore.metric + '</div> \n' +
        //         '<div class="metric-confidence">Confidence: ' + metricScore.confidence + '%</div>'
        // }
    }

    result += '</div></div>'
    return result
}

// subMetrics is list of tuples (nameOfMetric, MetricScore, MetricConf, barDisplay, highlight)
const getMetricContainerWSubContainers = (metricName, subMetrics) => {
    let subcontainers = ""

    for (let i = 0; i < subMetrics.length; i++) {
        subcontainers += '<div class="submetric-container">\n' + getMetricDisplay(subMetrics[i][1], subMetrics[i][0], subMetrics[i][2], false)  + '</div>'
    }

    return '<div class="metric-container">\n' +
        '<div class="metric-container-title">' + metricName + '</div>' + subcontainers + '</div>'
}

// activityScore, licenseScore, stars, contributors are tuples (metricScore, confidence, highlight)
const getBasicInfoDisplay = (owner, name, activityScore, licenseScore, stars, contributors) => {
    let result = '<div class="basic-info-display"> \n' +
    '<div class="basic-info" id="repoOwnerName">' + owner + '/' + name + '</div>'
    result += '<div class="basic-info-flexbox">'
   
    result += getMetricDisplay(activityScore, "Activity Score", false, true)
    result += getMetricDisplay(licenseScore, "License Score", false, true)
    result += getMetricDisplay(stars, "Stars", false, false)
    result += getMetricDisplay(contributors, "Contributors", false, false)

    // close div
    result += '</div></div>'
    return result
}

// Given an array of JSON objects, add highlight bool to highlight fields with max values
const AddHighlightJSON = (metricsAll) => {
    // Add highlight bool to each thing field of the JSON
    if (metricsAll.length < 1) {
        alert("Scores Unavailable")
    }

    for (var key in metricsAll[0]) {
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
        
        // owner, name, activityScore, licenseScore, stars, contributors
        result += getBasicInfoDisplay(metrics[i][0], metrics[i][1], metricsAll[i].repoActivityScore,
            metricsAll[i].repoLicenseScore, metricsAll[i].stars, metricsAll[i].contributors)

        result += '<div class="metrics-display">'
        result += '<div class="metric-category">Activity Score Breakdown</div>'
        result += getMetricDisplay(metricsAll[i].issueClosureTime, 'Issue Closure Time', true, false)
        result += getMetricDisplay(metricsAll[i].commitCadence, 'Commit Cadence', true, false)
        result += getMetricDisplay(metricsAll[i].releaseCadence, 'Release Cadence', true, false)
        result += getMetricDisplay(metricsAll[i].ageLastRelease, 'Age of Last Release', true, false)
        result += '</div >'

        result += '<div class="repo-dependency-score">'
        result += '<div class="metric-category">Dependency Scores</div>'
        result += getMetricDisplay(metricsAll[i].dependencyActivityScore, 'Dependency Activity Score', false, true)
        result += getMetricDisplay(metricsAll[i].dependencyLicenseScore, 'Dependency License Score', false, true)
        result += '</div >'


        result += '</div >'
    }

    return result
}

export default DisplayScores;
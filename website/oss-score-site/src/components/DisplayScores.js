import './DisplayScores.css';
import './Homepage.css';
import React from "react";
import { AiOutlineInfoCircle } from "react-icons/ai";
import ReactDOMServer from "react-dom/server";


// note: For stats where lower numbers are better, min and max are switched
let MetricStats = {

    "Activity Score-tooltip": "Overall activity score based on github metadata",

    "License-tooltip": "Direct mapping based on the license of the repo",

    "Stars-tooltip": "The number of stargazers for the repository",

    "Contributors-tooltip": "The number of unique users who have contributed to the project over the given query time frame",

    "Dependency Activity Score-tooltip": "Overall activity score based on github metadata for project dependencies",

    "Dependency License Score-tooltip": "Direct mapping based on the licenses of the dependencies",

    "Last Commit-tooltip": "Time since last commit",

    "PR Closure Time-min": 60,
    "PR Closure Time-max": 0,
    "PR Closure Time-units": "days",
    "PR Closure Time-tooltip": "Average time for a pull request in the project to be closed. NOTE: This is only calculated based on the closed PRs",

    "Issue Closure Time-min": 60,
    "Issue Closure Time-max": 0,
    "Issue Closure Time-units": "days",
    "Issue Closure Time-tooltip": "Average time for an issue in the project to be closed. NOTE: This is only calculated based on the closed issues",

    "Commit Cadence-min": 0,
    "Commit Cadence-max": 7,
    "Commit Cadence-units": "commits/week",
    "Commit Cadence-tooltip": "Average pace of commits in the project. Total number of commits divided by the query time frame",

    "Release Cadence-min": 0,
    "Release Cadence-max": 0.33,
    "Release Cadence-units": "releases/month",
    "Release Cadence-tooltip": "Average pace of releases in the project. Total number of releases divided by the query time frame",

    "Age of Last Release-min": 52,
    "Age of Last Release-max": 0,
    "Age of Last Release-units": "weeks",
    "Age of Last Release-tooltip": "Time since the last release release",
}


// TODO: Highlighting better metric, when lower is bette
const getMetricDisplay = (metricScore, metricName, barDisplay, outOfTen, lg) => {
    const infoLogo = ReactDOMServer.renderToStaticMarkup(<AiOutlineInfoCircle />);
    const infoLogoString = infoLogo.toString()
    
    // round metric 
    metricScore.metric = Math.round(metricScore.metric * 100) / 100
    metricScore.confidence = Math.round(metricScore.confidence)

    // shorten stars metric
    let starsOver1k = false
    if (metricName === "Stars") {
        if (metricScore.metric >= 1000) {
            metricScore.metric = Math.round(metricScore.metric / 1000)
            starsOver1k = true
        }
    }

    let result = ''

    if (barDisplay) {
        result += '<div class="metric-container bar-container" '
    } else if (lg) {
        result += '<div class="metric-container lg-container" '
    } else {
        result += '<div class="metric-container sm-container" '
    }

    if (metricScore.highlight) {
        result += 'style="background-color: #b0c4de;">'
    } else {
        result += 'style="background-color: #d3d3d3;">'
    }
    
    result += '<div class="tool-tip">' + infoLogoString +
        '<span class="tooltiptext">' + MetricStats[metricName + "-tooltip"] + '</span>\n' +
        '</div>'

    if (barDisplay) {
        result += '<div class="bar-title-metric">\n' +
                        '<div class="metric-container-title">' + metricName + '</div>\n' +
                        '<div class="metric-num">' + metricScore.metric + ' ' + MetricStats[metricName + "-units"] + '</div>\n' +
                    '</div>\n' +
                    '<div class="metrics">' 
    } else {
        result += '<div class="metric-container-title">' + metricName + '</div>\n' +
            '<div class="metrics">' 
    }

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

        result += '<div class="bar-and-conf">\n' +  
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
        if (metricName === 'License') {
            if (metricScore.license === '') {
                result += '<div class="metric-num"> N/A'
            } else {
                // make sure license name fits in the box
                if (metricScore.license.length > 12) {
                    result += '<div class="metric-num" style="font-size: 15px">' + metricScore.license
                } else if (metricScore.license.length > 10) {
                    result += '<div class="metric-num" style="font-size: 21px">' + metricScore.license
                }else {
                    result += '<div class="metric-num">' + metricScore.license
                }
                result += '<div class="metric-confidence">Score: ' + metricScore.metric + '/10</div>'
            }
        } else if (metricName === 'Last Commit') {
            let roundedMetric = Math.round(metricScore.metric)
            if (roundedMetric <= 0) {
                result += '<div class="metric-num"> <1 day'
            } else {
                result += '<div class="metric-num">' + roundedMetric
                if (roundedMetric === 1) {
                    result += ' day'
                } else {
                    result += ' days' 
                }
            }
        } else {
            result += '<div class="metric-num">' + metricScore.metric
        }

        if (outOfTen) {
            result += '/10'
        } else if (starsOver1k) {
            result += 'k'
        }

        result += '</div>'
        
        if (metricName !== 'License') {
            result += '<div class="metric-confidence">Confidence: ' + metricScore.confidence + '%</div>'
        }
    }

    result += '</div></div>'
    return result
}

// activityScore, licenseScore, stars, contributors are tuples (metricScore, confidence, highlight)
const getBasicInfoDisplay = (owner, name, activityScore, licenseInfo, stars, contributors, ageLastCommit) => {
    let result = '<div class="basic-info-display"> \n' +
        '<a class="basic-info" id="repoOwnerName" target="_blank" href = "https://github.com/' + owner + '/' + name + '">' + owner + '/' + name + '</a>'
    result += '<div class="info-grid-2-cols">'
   
    result += getMetricDisplay(activityScore, "Activity Score", false, true, true)
    result += getMetricDisplay(ageLastCommit, "Last Commit", false, false, true)
    result += '</div>'
    result += '<div class="info-grid-3-cols">'
    result += getMetricDisplay(contributors, "Contributors", false, false, false)
    result += getMetricDisplay(stars, "Stars", false, false, false)
    result += getMetricDisplay(licenseInfo, "License", false, false, false)
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
        if (key !== 'message' && key !== 'license') {
            let maxOfMetric = 0;
            let minOfMetric = Infinity
            let metricArray = [];
            // Find max of metrics and store values for highlighting
            for (let i = 0; i < metricsAll.length; i++) {
                metricArray.push(metricsAll[i][key].metric)
                if (key === "issueClosureTime" || key === "ageLastRelease" || key === "ageLastCommit") {
                    if (metricsAll[i][key].metric < minOfMetric) {
                        minOfMetric = metricsAll[i][key].metric;
                    }
                } else {
                    if (metricsAll[i][key].metric > maxOfMetric) {
                        maxOfMetric = metricsAll[i][key].metric;
                    }
                }
            }

            // Add highligt property metrics
            for (let i = 0; i < metricArray.length; i++) {
                if (key === "issueClosureTime" || key === "ageLastRelease" || key === "ageLastCommit") {
                    if (metricArray[i] === minOfMetric) {
                        // Highlight metric
                        metricsAll[i][key].highlight = true;
                    } else {
                        // Don't highlight
                        metricsAll[i][key].highlight = false;
                    }
                } else {
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
        let licenseInfo = {
            'license': metricsAll[i].license,
            'metric': metricsAll[i].repoLicenseScore.metric,
            'confidence': metricsAll[i].repoLicenseScore.confidence,
            'highlight': metricsAll[i].repoLicenseScore.highlight
        } 
        // owner, name, activityScore, licenseScore, stars, contributors
        result += getBasicInfoDisplay(metrics[i][0], metrics[i][1], metricsAll[i].repoActivityScore,
            licenseInfo, metricsAll[i].stars, metricsAll[i].contributors, metricsAll[i].ageLastCommit)

        result += '<div class="metrics-display">'
        result += '<div class="metric-category">Activity Score Breakdown</div>'
        result += getMetricDisplay(metricsAll[i].issueClosureTime, 'Issue Closure Time', true, false, false)
        result += getMetricDisplay(metricsAll[i].commitCadence, 'Commit Cadence', true, false, false)
        result += getMetricDisplay(metricsAll[i].prClosureTime, 'PR Closure Time', true, false, false)
        result += getMetricDisplay(metricsAll[i].releaseCadence, 'Release Cadence', true, false, false)
        result += getMetricDisplay(metricsAll[i].ageLastRelease, 'Age of Last Release', true, false, false)
        result += '</div >'

        result += '<div class="repo-dependency-score">'
        result += '<div class="metric-category">Dependency Scores</div>'
        result += '<div class="info-grid-2-cols">'
        result += getMetricDisplay(metricsAll[i].dependencyActivityScore, 'Dependency Activity Score', false, true, true)
        result += getMetricDisplay(metricsAll[i].dependencyLicenseScore, 'Dependency License Score', false, true, true)
        result += '</div >'
        result += '</div >'


        result += '</div >'
    }

    return result
}

export default DisplayScores;
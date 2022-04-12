const basePath = 'https://hvacjx4u1l.execute-api.us-east-2.amazonaws.com/prod/catalog/github' //prod
//const basePath = 'https://xvzhkajkzh.execute-api.us-east-2.amazonaws.com/dev/catalog/github' //dev
const ossScoreSite = 'https://oss-score.herokuapp.com'
const calculationMessages = ['Score not yet calculated', 'Error querying score', 'Data out of date']
const errorMessages = ['Could not access repo, check that it was inputted correctly and is public', 'Cannot provide score for private repo']

function promiseTimeout (time) {
    return new Promise(function(resolve, reject) {
        setTimeout(function() {resolve();}, time);
    });
};

async function requestScores(owner, repo) {
    let message = null;
    let success = false
    let requestURL = basePath + '/owner/' + owner + '/name/' + repo;
    let requestBody = null;
    if (TimeFrame != null){
        requestBody = JSON.stringify({"timeFrame": parseInt(TimeFrame)});
    }
    let promise = 
        fetch(requestURL, {
            method: 'POST',
            mode: 'cors',
            body: requestBody
        }).then(async (response) => {
            if (response.status == 201) {
                let messagePromise = response.json();
                await messagePromise.then(response => {
                    message = response.message;
                    success = true
                }).catch(err => {
                    console.error(err);
                });
            } else if (response.status == 406)  {
                message = "Cannot provide score for private repo";
            } else if ((response.status == 501) || (response.status == 503))  {
                message = "Error: Internal Servor Error";
            } else {
                message = "Un-handled response from OSS-Score API";
            }
        }).catch(err => {
            message = "Error while placing/parsing post";
            message += err;
            console.error(err);
        });
    await promise;

    return {
        message: message,
        success: success,
    };
}

async function updateScores(scoreDiv, owner, repo) {
    promiseTimeout(15000).then(() => {
        console.log('Updating Score');
        getScores(owner, repo).then(scores => {
            if ((scores.activity != null) && (scores.license != null)) {
                insertScores(scoreDiv, scores);
                if (scores.depRatio < 0.9) {
                    updateScores(scoreDiv, owner, repo);
                }
            } else {
                console.log("Error retrieving scores in updateScores");
            }
        });
    });
}

async function awaitResults(scoreDiv, owner, repo) {
    promiseTimeout(1000).then(() => {
        console.log('Requesting Score');
        getScores(owner, repo).then(scores => {
            if (scores.activity != null && scores.license != null) {
                insertScores(scoreDiv, scores);
                updateScores(scoreDiv, owner, repo);
            } else {
                awaitResults(scoreDiv, owner, repo);
            }
        });
    });
}

function hslToRgb(h, s, l){
    var r, g, b;

    if(s == 0){
        r = g = b = l; // achromatic
    }else{
        function hue2rgb(p, q, t){
            if(t < 0) t += 1;
            if(t > 1) t -= 1;
            if(t < 1/6) return p + (q - p) * 6 * t;
            if(t < 1/2) return q;
            if(t < 2/3) return p + (q - p) * (2/3 - t) * 6;
            return p;
        }

        var q = l < 0.5 ? l * (1 + s) : l + s - l * s;
        var p = 2 * l - q;
        r = hue2rgb(p, q, h + 1/3);
        g = hue2rgb(p, q, h);
        b = hue2rgb(p, q, h - 1/3);
    }

    return [Math.floor(r * 255), Math.floor(g * 255), Math.floor(b * 255)];
}

function numberToColorHsl(i, min, max) {
    var ratio = i;
    if (min> 0 || max < 1) {
        if (i < min) {
            ratio = 0;
        } else if (i > max) {
            ratio = 1;
        } else {
            var range = max - min;
            ratio = (i-min) / range;
        }
    }
    
    // as the function expects a value between 0 and 1, and red = 0° and green = 120°
    // we convert the input to the appropriate hue value
    var hue = ratio * 1.2 / 3.60;
    //if (minMaxFactor!=1) hue /= minMaxFactor;
    //console.log(hue);
    
    // we convert hsl to rgb (saturation 100%, lightness 50%)
    var rgb = hslToRgb(hue, 1, .5);
    // we format to css value and return
    return 'rgb(' + rgb[0] + ',' + rgb[1] + ',' + rgb[2] + ')'; 
}

// insertHTML expects the div to insert the html into, the scores to be populated (may be null), and the 
// state of the extension. Depending on the state it inserts the correct html and then returns nothing. 
function insertHTML(scoreDiv, scores, state) {
    let header = "<h2 class=\"h4 mb-3\"> <a style='text-decoration: none; color: rgb(0,0,0)' href=" + ossScoreSite + "/generate-scores>OSS Score</a> </h2>";
    let loading_gears = "<div class='loading-extension'>     <svg class='machine-extension-light'xmlns='http://www.w3.org/2000/svg' x='0px' y='0px' viewBox='0 0 645 526'>       <defs/>       <g>         <path  x='-173,694' y='-173,694' class='large-shadow-extension-light' d='M645 194v-21l-29-4c-1-10-3-19-6-28l25-14 -8-19 -28 7c-5-8-10-16-16-24L602 68l-15-15 -23 17c-7-6-15-11-24-16l7-28 -19-8 -14 25c-9-3-18-5-28-6L482 10h-21l-4 29c-10 1-19 3-28 6l-14-25 -19 8 7 28c-8 5-16 10-24 16l-23-17L341 68l17 23c-6 7-11 15-16 24l-28-7 -8 19 25 14c-3 9-5 18-6 28l-29 4v21l29 4c1 10 3 19 6 28l-25 14 8 19 28-7c5 8 10 16 16 24l-17 23 15 15 23-17c7 6 15 11 24 16l-7 28 19 8 14-25c9 3 18 5 28 6l4 29h21l4-29c10-1 19-3 28-6l14 25 19-8 -7-28c8-5 16-10 24-16l23 17 15-15 -17-23c6-7 11-15 16-24l28 7 8-19 -25-14c3-9 5-18 6-28L645 194zM471 294c-61 0-110-49-110-110S411 74 471 74s110 49 110 110S532 294 471 294z'/>       </g>       <g>         <path x='-136,996' y='-136,996' class='medium-shadow-extension-light' d='M402 400v-21l-28-4c-1-10-4-19-7-28l23-17 -11-18L352 323c-6-8-13-14-20-20l11-26 -18-11 -17 23c-9-4-18-6-28-7l-4-28h-21l-4 28c-10 1-19 4-28 7l-17-23 -18 11 11 26c-8 6-14 13-20 20l-26-11 -11 18 23 17c-4 9-6 18-7 28l-28 4v21l28 4c1 10 4 19 7 28l-23 17 11 18 26-11c6 8 13 14 20 20l-11 26 18 11 17-23c9 4 18 6 28 7l4 28h21l4-28c10-1 19-4 28-7l17 23 18-11 -11-26c8-6 14-13 20-20l26 11 11-18 -23-17c4-9 6-18 7-28L402 400zM265 463c-41 0-74-33-74-74 0-41 33-74 74-74 41 0 74 33 74 74C338 430 305 463 265 463z'/>       </g>       <g >         <path x='-100,136' y='-100,136' class='small-shadow-extension-light' d='M210 246v-21l-29-4c-2-10-6-18-11-26l18-23 -15-15 -23 18c-8-5-17-9-26-11l-4-29H100l-4 29c-10 2-18 6-26 11l-23-18 -15 15 18 23c-5 8-9 17-11 26L10 225v21l29 4c2 10 6 18 11 26l-18 23 15 15 23-18c8 5 17 9 26 11l4 29h21l4-29c10-2 18-6 26-11l23 18 15-15 -18-23c5-8 9-17 11-26L210 246zM110 272c-20 0-37-17-37-37s17-37 37-37c20 0 37 17 37 37S131 272 110 272z'/>       </g>       <g>         <path x='-100,136' y='-100,136' class='small-extension' d='M200 236v-21l-29-4c-2-10-6-18-11-26l18-23 -15-15 -23 18c-8-5-17-9-26-11l-4-29H90l-4 29c-10 2-18 6-26 11l-23-18 -15 15 18 23c-5 8-9 17-11 26L0 215v21l29 4c2 10 6 18 11 26l-18 23 15 15 23-18c8 5 17 9 26 11l4 29h21l4-29c10-2 18-6 26-11l23 18 15-15 -18-23c5-8 9-17 11-26L200 236zM100 262c-20 0-37-17-37-37s17-37 37-37c20 0 37 17 37 37S121 262 100 262z'/>       </g>       <g>         <path x='-173,694' y='-173,694' class='large-extension' d='M635 184v-21l-29-4c-1-10-3-19-6-28l25-14 -8-19 -28 7c-5-8-10-16-16-24L592 58l-15-15 -23 17c-7-6-15-11-24-16l7-28 -19-8 -14 25c-9-3-18-5-28-6L472 0h-21l-4 29c-10 1-19 3-28 6L405 9l-19 8 7 28c-8 5-16 10-24 16l-23-17L331 58l17 23c-6 7-11 15-16 24l-28-7 -8 19 25 14c-3 9-5 18-6 28l-29 4v21l29 4c1 10 3 19 6 28l-25 14 8 19 28-7c5 8 10 16 16 24l-17 23 15 15 23-17c7 6 15 11 24 16l-7 28 19 8 14-25c9 3 18 5 28 6l4 29h21l4-29c10-1 19-3 28-6l14 25 19-8 -7-28c8-5 16-10 24-16l23 17 15-15 -17-23c6-7 11-15 16-24l28 7 8-19 -25-14c3-9 5-18 6-28L635 184zM461 284c-61 0-110-49-110-110S401 64 461 64s110 49 110 110S522 284 461 284z'/>       </g>       <g>         <path x='-136,996' y='-136,996' class='medium-extension' d='M392 390v-21l-28-4c-1-10-4-19-7-28l23-17 -11-18L342 313c-6-8-13-14-20-20l11-26 -18-11 -17 23c-9-4-18-6-28-7l-4-28h-21l-4 28c-10 1-19 4-28 7l-17-23 -18 11 11 26c-8 6-14 13-20 20l-26-11 -11 18 23 17c-4 9-6 18-7 28l-28 4v21l28 4c1 10 4 19 7 28l-23 17 11 18 26-11c6 8 13 14 20 20l-11 26 18 11 17-23c9 4 18 6 28 7l4 28h21l4-28c10-1 19-4 28-7l17 23 18-11 -11-26c8-6 14-13 20-20l26 11 11-18 -23-17c4-9 6-18 7-28L392 390zM255 453c-41 0-74-33-74-74 0-41 33-74 74-74 41 0 74 33 74 74C328 420 295 453 255 453z'/>       </g>     </svg> </div>";
    let isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    let button = '<button class=requestScore-light id=requestButton> Request Score </button>'
    if (isDark) {
        scoreDiv.style.backgroundColor = "rgb(13, 17, 23)";
        scoreDiv.style.color = "rgb(201,209,217)";
        header = "<h2 class=\"h4 mb-3\"> <a style='text-decoration: none; color: rgb(201,209,217)' href=" + ossScoreSite + "/generate-scores>OSS Score</a> </h2>";
        loading_gears = "<div class='loading-extension'>     <svg class='machine-extension-dark'xmlns='http://www.w3.org/2000/svg' x='0px' y='0px' viewBox='0 0 645 526'>       <defs/>       <g>         <path  x='-173,694' y='-173,694' class='large-shadow-extension-dark' d='M645 194v-21l-29-4c-1-10-3-19-6-28l25-14 -8-19 -28 7c-5-8-10-16-16-24L602 68l-15-15 -23 17c-7-6-15-11-24-16l7-28 -19-8 -14 25c-9-3-18-5-28-6L482 10h-21l-4 29c-10 1-19 3-28 6l-14-25 -19 8 7 28c-8 5-16 10-24 16l-23-17L341 68l17 23c-6 7-11 15-16 24l-28-7 -8 19 25 14c-3 9-5 18-6 28l-29 4v21l29 4c1 10 3 19 6 28l-25 14 8 19 28-7c5 8 10 16 16 24l-17 23 15 15 23-17c7 6 15 11 24 16l-7 28 19 8 14-25c9 3 18 5 28 6l4 29h21l4-29c10-1 19-3 28-6l14 25 19-8 -7-28c8-5 16-10 24-16l23 17 15-15 -17-23c6-7 11-15 16-24l28 7 8-19 -25-14c3-9 5-18 6-28L645 194zM471 294c-61 0-110-49-110-110S411 74 471 74s110 49 110 110S532 294 471 294z'/>       </g>       <g>         <path x='-136,996' y='-136,996' class='medium-shadow-extension-dark' d='M402 400v-21l-28-4c-1-10-4-19-7-28l23-17 -11-18L352 323c-6-8-13-14-20-20l11-26 -18-11 -17 23c-9-4-18-6-28-7l-4-28h-21l-4 28c-10 1-19 4-28 7l-17-23 -18 11 11 26c-8 6-14 13-20 20l-26-11 -11 18 23 17c-4 9-6 18-7 28l-28 4v21l28 4c1 10 4 19 7 28l-23 17 11 18 26-11c6 8 13 14 20 20l-11 26 18 11 17-23c9 4 18 6 28 7l4 28h21l4-28c10-1 19-4 28-7l17 23 18-11 -11-26c8-6 14-13 20-20l26 11 11-18 -23-17c4-9 6-18 7-28L402 400zM265 463c-41 0-74-33-74-74 0-41 33-74 74-74 41 0 74 33 74 74C338 430 305 463 265 463z'/>       </g>       <g >         <path x='-100,136' y='-100,136' class='small-shadow-extension-dark' d='M210 246v-21l-29-4c-2-10-6-18-11-26l18-23 -15-15 -23 18c-8-5-17-9-26-11l-4-29H100l-4 29c-10 2-18 6-26 11l-23-18 -15 15 18 23c-5 8-9 17-11 26L10 225v21l29 4c2 10 6 18 11 26l-18 23 15 15 23-18c8 5 17 9 26 11l4 29h21l4-29c10-2 18-6 26-11l23 18 15-15 -18-23c5-8 9-17 11-26L210 246zM110 272c-20 0-37-17-37-37s17-37 37-37c20 0 37 17 37 37S131 272 110 272z'/>       </g>       <g>         <path x='-100,136' y='-100,136' class='small-extension' d='M200 236v-21l-29-4c-2-10-6-18-11-26l18-23 -15-15 -23 18c-8-5-17-9-26-11l-4-29H90l-4 29c-10 2-18 6-26 11l-23-18 -15 15 18 23c-5 8-9 17-11 26L0 215v21l29 4c2 10 6 18 11 26l-18 23 15 15 23-18c8 5 17 9 26 11l4 29h21l4-29c10-2 18-6 26-11l23 18 15-15 -18-23c5-8 9-17 11-26L200 236zM100 262c-20 0-37-17-37-37s17-37 37-37c20 0 37 17 37 37S121 262 100 262z'/>       </g>       <g>         <path x='-173,694' y='-173,694' class='large-extension' d='M635 184v-21l-29-4c-1-10-3-19-6-28l25-14 -8-19 -28 7c-5-8-10-16-16-24L592 58l-15-15 -23 17c-7-6-15-11-24-16l7-28 -19-8 -14 25c-9-3-18-5-28-6L472 0h-21l-4 29c-10 1-19 3-28 6L405 9l-19 8 7 28c-8 5-16 10-24 16l-23-17L331 58l17 23c-6 7-11 15-16 24l-28-7 -8 19 25 14c-3 9-5 18-6 28l-29 4v21l29 4c1 10 3 19 6 28l-25 14 8 19 28-7c5 8 10 16 16 24l-17 23 15 15 23-17c7 6 15 11 24 16l-7 28 19 8 14-25c9 3 18 5 28 6l4 29h21l4-29c10-1 19-3 28-6l14 25 19-8 -7-28c8-5 16-10 24-16l23 17 15-15 -17-23c6-7 11-15 16-24l28 7 8-19 -25-14c3-9 5-18 6-28L635 184zM461 284c-61 0-110-49-110-110S401 64 461 64s110 49 110 110S522 284 461 284z'/>       </g>       <g>         <path x='-136,996' y='-136,996' class='medium-extension' d='M392 390v-21l-28-4c-1-10-4-19-7-28l23-17 -11-18L342 313c-6-8-13-14-20-20l11-26 -18-11 -17 23c-9-4-18-6-28-7l-4-28h-21l-4 28c-10 1-19 4-28 7l-17-23 -18 11 11 26c-8 6-14 13-20 20l-26-11 -11 18 23 17c-4 9-6 18-7 28l-28 4v21l28 4c1 10 4 19 7 28l-23 17 11 18 26-11c6 8 13 14 20 20l-11 26 18 11 17-23c9 4 18 6 28 7l4 28h21l4-28c10-1 19-4 28-7l17 23 18-11 -11-26c8-6 14-13 20-20l26 11 11-18 -23-17c4-9 6-18 7-28L392 390zM255 453c-41 0-74-33-74-74 0-41 33-74 74-74 41 0 74 33 74 74C328 420 295 453 255 453z'/>       </g>     </svg> </div>";
        button = '<button class=requestScore-dark id=requestButton> Request Score </button>';
    }
    scoreDiv.innerHTML = header;
    if (state == "ready") {
        let activityColor = numberToColorHsl(scores.activity.score/10, 0, 1);
        scoreDiv.innerHTML += '<dot-extension style="color:' + activityColor + '">●</dot-extension>';
        scoreDiv.innerHTML += 'Activity: ' + (scores.activity.score).toFixed(1) + ' of 10';
        scoreDiv.innerHTML += '&nbsp; | &nbsp; Confidence: ' + scores.activity.confidence.toFixed(0) + '%';
        scoreDiv.innerHTML += '<br/><br/>';
        let licenseColor = numberToColorHsl(scores.license.score/10, 0, 1);
        scoreDiv.innerHTML += '<dot-extension style="color:' + licenseColor + '">●</dot-extension>';
        scoreDiv.innerHTML += 'License: ' + (scores.license.score).toFixed(1) + ' of 10';
        scoreDiv.innerHTML += '&nbsp; | &nbsp; Confidence: ' + scores.license.confidence.toFixed(0) + '%';
    } else if (state == "loading") {
        scoreDiv.innerHTML += loading_gears;
    } else if (state == "uncalculated") {
        scoreDiv.innerHTML += scores.message;
        scoreDiv.innerHTML += '<br><br>';
        scoreDiv.innerHTML += button;
    } else if (state == "error") {
        scoreDiv.innerHTML += scores.message;
    }
}

function insertScores(scoreDiv, scores) {
    insertHTML(scoreDiv, scores, "ready");
}

async function insertScoreSection(owner, repo, scoreDiv, scoresPromise) {
    //inject into correct part of site
    try {
        let repoInfo = document.querySelectorAll('.BorderGrid-row');
        let releases = repoInfo[1];
        let parent = releases.parentNode;
        parent.insertBefore(scoreDiv, releases);
    } catch (error) {
        console.log("Error in insertScoreSection: " + error);
        return;
    }
    insertHTML(scoreDiv, null, "loading");
    scoresPromise.then(scores => {
        if (scores.activity != null && scores.license != null) { // VALID SCORES RETURNED
            insertScores(scoreDiv, scores);
            updateScores(scoreDiv, owner, repo);
        } else if (calculationMessages.includes(scores.message)) { // SCORES NEED TO BE CALCULATED
            insertHTML(scoreDiv, scores, "uncalculated");
            document.getElementById('requestButton').addEventListener('click', async function() {
                requestScores(owner, repo).then(requestResponse => {
                    if (requestResponse.success) {
                        insertHTML(scoreDiv, null, "loading");
                        awaitResults(scoreDiv, owner, repo);
                    } else {
                        insertHTML(scoreDiv, requestResponse, "error");
                    }
                });
            });
        } else if (errorMessages.includes(scores.message)) { // SCORES CAN NOT BE CALCULATED
            insertHTML(scoreDiv, scores, "error");
        } else { // SCORE CALCULATION HAPPENING
            insertHTML(scoreDiv, null, "loading");
        }
    });
}

async function getScores(owner, repo) {
    let scores = {license: null, activity: null, message: null, depRatio: 0};
    let promises = [];
    let queryParams = '';    
    if (TimeFrame != null) {
        queryParams = "?timeFrame=" + TimeFrame;
    }
    let licenseRequestUrl = basePath + '/owner/' + owner + '/name/' + repo + '/type/license' + queryParams;
    promises.push(
        fetch(licenseRequestUrl).then(async (response) => {
            if (response.status == 200) {
                let scorePromise = response.json();
                await scorePromise.then(response => {
                    if (response.message == "Score ready") {
                        scores.license = response.score;
                        scores.depRatio += response.depRatio / 2;
                    } else {
                        scores.message = response.message;
                    };
                }).catch(err => {
                    console.error(err);
                });
            } else if (response.status == 406)  {
                scores.message = "Cannot provide score for private repo";
            } else if ((response.status == 501) || (response.status == 503))  {
                scores.message = "Error: Internal Servor Error";
            } else {
                scores.message = "Un-handled response from OSS-Score API";
            }
        }).catch(err => {
            console.error(err);
        })
    );

    let activityRequestUrl = basePath + '/owner/' + owner + '/name/' + repo + '/type/activity' + queryParams;
    promises.push(
        fetch(activityRequestUrl).then(async (response) => {
            if (response.status == 200) {
                let scorePromise = response.json();
                await scorePromise.then(response => {
                    if (response.message == "Score ready") {
                        scores.activity = response.score;
                        scores.depRatio += response.depRatio / 2;
                    } else {
                        scores.message = response.message;
                    };
                }).catch(err => {
                    console.error(err);
                });
            } else if (response.status == 406)  {
                scores.message = "Cannot provide score for private repo";
            } else if ((response.status == 501) || (response.status == 503))  {
                scores.message = "Error: cannot calculate score request";
            } else {
                scores.message = "Un-handled response from OSS-Score API";
            }
        }).catch(err => {
            console.error(err);
        })
    );

    await Promise.all(promises);

    return scores;
}

let url = window.location.href.replace('https://github.com/', '');

let splitUrl = url.split('/').filter(element => element != '');
let owner = '';
let repo = '';

if (splitUrl.length == 2) { // Repo homepage
    owner = splitUrl[0];
    repo = splitUrl[1];
} else if (splitUrl.length == 4) {
    if (splitUrl[2] == 'tree') { // on specific branch
        owner = splitUrl[0];
        repo = splitUrl[1];
    }
}

var TimeFrame = null;
chrome.storage.sync.get(['key'], function(result) {
    console.log('Value currently is ' + result.key);
    TimeFrame = result.key
});

if (owner != '' && repo != '') {
    let scoreDiv = document.createElement('div');
    scoreDiv.className = 'BorderGrid-cell';
    
    chrome.runtime.onMessage.addListener(
        function(request, sender, sendResponse) {
          TimeFrame = request.timeFrame;
          sendResponse({message: "recieved new time frame"});
          insertScoreSection(owner, repo, scoreDiv, getScores(owner, repo, scoreDiv));
        }
      );
    insertScoreSection(owner, repo, scoreDiv, getScores(owner, repo, scoreDiv));
}
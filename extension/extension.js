const basePath = 'https://oss-hub.herokuapp.com/'
//const basePath = ''

function promiseTimeout (time) {
    return new Promise(function(resolve, reject) {
        setTimeout(function() {resolve();}, time);
    });
};

async function requestScoreCalculation(path, scoreType) {
    fetch(basePath + path + '/' + scoreType + '/score', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        mode: 'cors',
    }).then(response=> {
        switch (response.status) {
            case 201:
                return 'Request Recieved \n Calculating Scores';
            case 404:
                return 'Error while requesting scores';
            default:
                return 'Error: ' + response.status;
        }
    }).catch(err => {
        console.error(err);
        return 'Internal Error';
    })
}

function insertScores(scoreDiv, scores) {
    scoreDiv.innerHTML = "<h2 class=\"h4 mb-3\"> OSS Scores </h2>"

    var colorString1 = "good_score";
    var colorString2 = "good_score";
    if (scores.license.score <= 50) {
        colorString1 = "bad_score";
    }
    if (scores.activity.score <= 50) {
        colorString2 = "bad_score";
    }
    scoreDiv.innerHTML += "<span class='" + colorString1 + "'>" + 'License: ' + scores.license.score + "</span>";
    scoreDiv.innerHTML += '&nbsp; | &nbsp; Confidence: ' + scores.license.confidence + '%';
    scoreDiv.innerHTML += '<br/><br/>'
    scoreDiv.innerHTML += "<span class='" + colorString2 + "'>" + 'Activity: ' + scores.activity.score + "</span>";
    scoreDiv.innerHTML += '&nbsp; | &nbsp; Confidence: ' + scores.activity.confidence + '%';
}

async function insertScoreSection(path, scoreDiv, scoresPromise) {
    
    //inject into correct part of site
    let repoInfo = document.querySelectorAll('.BorderGrid-row');
    let releases = repoInfo[1];
    let parent = releases.parentNode;
    parent.insertBefore(scoreDiv, releases);

    scoreDiv.innerHTML = "<h2 class=\"h4 mb-3\"> OSS Scores </h2> Scores Loading..."

    // let scores = await scoresPromise;
    scoresPromise.then(scores => {
        if (scores.message == null || scores.message == '') { // VALID SCORES RETURNED
            insertScores(scoreDiv, scores);
        } else { // NO SCORES
            scoreDiv.innerHTML = "<h2 class=\"h4 mb-3\"> OSS Scores </h2>"
            scoreDiv.innerHTML += scores.message;


            // add conditional to replace if message is waiting and add loading
            if (scores.message.includes('not yet calculated') || scores.message.includes('in progress')) {
                scoreDiv.innerHTML = "<h2 class=\"h4 mb-3\"> OSS Scores </h2>";
                let image_url = chrome.runtime.getURL("images/loading-spinning.gif");
                //scoreDiv.src = image_url;
                //scoreDiv.appendChild(document.createElement('img')).src = image_url;
                scoreDiv.innerHTML += "<svg width='100' height='100'> <circle cx='50' cy='50' r='40' stroke='green' stroke-width='4' fill='yellow' /> </svg>";
                //scoreDiv.innerHTML += "<img src='" + image_url + "' alt='loading' width='100%'></img>"
                //awaitResults(scoreDiv, path);
            }

            if (scores.message.includes('cached')) {
                scoreDiv.innerHTML += '<br><br>'
                scoreDiv.innerHTML += '<button class=requestScore id=requestButton> Request Score </button>'

                document.getElementById('requestButton').addEventListener('click', async function() {
                    scoreDiv.innerHTML = "<h2 class=\"h4 mb-3\"> OSS Scores </h2>"

                    //maybe put this in try catch?
                    requestScoreCalculation(path, 'activity');
                    requestScoreCalculation(path, 'license');
                    scoreDiv.innerHTML += 'Scores Requested';
                    awaitResults(scoreDiv, path);
                });
            }
        }
    })

}

async function awaitResults(scoreDiv, repoPath) {
    promiseTimeout(2500).then(() => {
        console.log('Requesting Score');
        getScores(repoPath).then(scores => {
            message = scores.message;

            if (message == undefined|| message == '') {
                insertScores(scoreDiv, scores);
            } else {
                awaitResults(scoreDiv, repoPath);
            }
        });
    });
}

async function getFakeScores(repoPath) {
    let scores = {license: {score: null, confidence: null}, activity: {score: null, confidence: null}, message: null};
    let score1 = 50;
    let score2 = 88;
    scores.license.score = score1;
    scores.activity.score = score2;
    scores.license.confidence = 100;
    scores.activity.confidence = 100;
    scores.message = "not yet calculated";
    return scores;
}

async function getScores(repoPath) {
    let scores = {license: null, activity: null, message: null};
    let promises = [];    
    let licenseRequestUrl = basePath + repoPath + '/license/score';
    promises.push(
        fetch(licenseRequestUrl).then(async (response) => {
            if (response.status == 200) {
                let scorePromise = response.json();
                await scorePromise.then(score => {
                    scores.license = score;
                }).catch(err => {
                    console.error(err);
                });
            } else {
                let messagePromise = response.json();
                await messagePromise.then(message => {
                    scores.message = message.message;
                }).catch(err => {
                    console.error(err);
                });
        }
        }).catch(err => {
            console.error(err);
        })
    );

    let activityRequestUrl = basePath + repoPath + '/activity/score';
    promises.push(
        fetch(activityRequestUrl).then(async (response) => {
            if (response.status == 200) {
                let scorePromise = response.json();
                await scorePromise.then(score => {
                    scores.activity = score;
                }).catch(err => {
                    console.error(err);
                });
            } else {
                let messagePromise = response.json();
                await messagePromise.then(message => {
                    scores.message = message.message;
                }).catch(err => {
                    console.error(err);
                });
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

if (owner != '' && repo != '') {
    //let rowDiv= document.createElement('div');
    //rowDiv.className = 'BorderGrid-row';
 
    let scoreDiv = document.createElement('div');
    scoreDiv.className = 'BorderGrid-cell';

    let path = 'github/' + owner + '/' + repo;
    //insertScoreSection(path, scoreDiv, getScores(path, scoreDiv));
    insertScoreSection(path, scoreDiv, getFakeScores(path, scoreDiv));
}
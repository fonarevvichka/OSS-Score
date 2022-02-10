const basePath = 'https://ru8ibij7yc.execute-api.us-east-2.amazonaws.com/staging/catalog/github'

function promiseTimeout (time) {
    return new Promise(function(resolve, reject) {
        setTimeout(function() {resolve();}, time);
    });
};

async function requestScoreCalculation(owner, repo, scoreType) {

    fetch(basePath + '/owner/' + owner + '/repo/' + repo + '/type/' + scoreType, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        mode: 'cors',
    }).then(response=> {
        switch (response.status) {
            case 200:
                return 'Request Recieved \n Calculating Scores';
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

async function insertScoreSection(owner, repo, scoreDiv, scoresPromise) {
    
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
                scoreDiv.innerHTML += 'Scores Loading, Please Wait'
                scoreDiv.innerHTML += '<br><br>'
                let loading_gears = "<div class='loading'>     <svg class='machine'xmlns='http://www.w3.org/2000/svg' x='0px' y='0px' viewBox='0 0 645 526'>       <defs/>       <g>         <path  x='-173,694' y='-173,694' class='large-shadow' d='M645 194v-21l-29-4c-1-10-3-19-6-28l25-14 -8-19 -28 7c-5-8-10-16-16-24L602 68l-15-15 -23 17c-7-6-15-11-24-16l7-28 -19-8 -14 25c-9-3-18-5-28-6L482 10h-21l-4 29c-10 1-19 3-28 6l-14-25 -19 8 7 28c-8 5-16 10-24 16l-23-17L341 68l17 23c-6 7-11 15-16 24l-28-7 -8 19 25 14c-3 9-5 18-6 28l-29 4v21l29 4c1 10 3 19 6 28l-25 14 8 19 28-7c5 8 10 16 16 24l-17 23 15 15 23-17c7 6 15 11 24 16l-7 28 19 8 14-25c9 3 18 5 28 6l4 29h21l4-29c10-1 19-3 28-6l14 25 19-8 -7-28c8-5 16-10 24-16l23 17 15-15 -17-23c6-7 11-15 16-24l28 7 8-19 -25-14c3-9 5-18 6-28L645 194zM471 294c-61 0-110-49-110-110S411 74 471 74s110 49 110 110S532 294 471 294z'/>       </g>       <g>         <path x='-136,996' y='-136,996' class='medium-shadow' d='M402 400v-21l-28-4c-1-10-4-19-7-28l23-17 -11-18L352 323c-6-8-13-14-20-20l11-26 -18-11 -17 23c-9-4-18-6-28-7l-4-28h-21l-4 28c-10 1-19 4-28 7l-17-23 -18 11 11 26c-8 6-14 13-20 20l-26-11 -11 18 23 17c-4 9-6 18-7 28l-28 4v21l28 4c1 10 4 19 7 28l-23 17 11 18 26-11c6 8 13 14 20 20l-11 26 18 11 17-23c9 4 18 6 28 7l4 28h21l4-28c10-1 19-4 28-7l17 23 18-11 -11-26c8-6 14-13 20-20l26 11 11-18 -23-17c4-9 6-18 7-28L402 400zM265 463c-41 0-74-33-74-74 0-41 33-74 74-74 41 0 74 33 74 74C338 430 305 463 265 463z'/>       </g>       <g >         <path x='-100,136' y='-100,136' class='small-shadow' d='M210 246v-21l-29-4c-2-10-6-18-11-26l18-23 -15-15 -23 18c-8-5-17-9-26-11l-4-29H100l-4 29c-10 2-18 6-26 11l-23-18 -15 15 18 23c-5 8-9 17-11 26L10 225v21l29 4c2 10 6 18 11 26l-18 23 15 15 23-18c8 5 17 9 26 11l4 29h21l4-29c10-2 18-6 26-11l23 18 15-15 -18-23c5-8 9-17 11-26L210 246zM110 272c-20 0-37-17-37-37s17-37 37-37c20 0 37 17 37 37S131 272 110 272z'/>       </g>       <g>         <path x='-100,136' y='-100,136' class='small' d='M200 236v-21l-29-4c-2-10-6-18-11-26l18-23 -15-15 -23 18c-8-5-17-9-26-11l-4-29H90l-4 29c-10 2-18 6-26 11l-23-18 -15 15 18 23c-5 8-9 17-11 26L0 215v21l29 4c2 10 6 18 11 26l-18 23 15 15 23-18c8 5 17 9 26 11l4 29h21l4-29c10-2 18-6 26-11l23 18 15-15 -18-23c5-8 9-17 11-26L200 236zM100 262c-20 0-37-17-37-37s17-37 37-37c20 0 37 17 37 37S121 262 100 262z'/>       </g>       <g>         <path x='-173,694' y='-173,694' class='large' d='M635 184v-21l-29-4c-1-10-3-19-6-28l25-14 -8-19 -28 7c-5-8-10-16-16-24L592 58l-15-15 -23 17c-7-6-15-11-24-16l7-28 -19-8 -14 25c-9-3-18-5-28-6L472 0h-21l-4 29c-10 1-19 3-28 6L405 9l-19 8 7 28c-8 5-16 10-24 16l-23-17L331 58l17 23c-6 7-11 15-16 24l-28-7 -8 19 25 14c-3 9-5 18-6 28l-29 4v21l29 4c1 10 3 19 6 28l-25 14 8 19 28-7c5 8 10 16 16 24l-17 23 15 15 23-17c7 6 15 11 24 16l-7 28 19 8 14-25c9 3 18 5 28 6l4 29h21l4-29c10-1 19-3 28-6l14 25 19-8 -7-28c8-5 16-10 24-16l23 17 15-15 -17-23c6-7 11-15 16-24l28 7 8-19 -25-14c3-9 5-18 6-28L635 184zM461 284c-61 0-110-49-110-110S401 64 461 64s110 49 110 110S522 284 461 284z'/>       </g>       <g>         <path x='-136,996' y='-136,996' class='medium' d='M392 390v-21l-28-4c-1-10-4-19-7-28l23-17 -11-18L342 313c-6-8-13-14-20-20l11-26 -18-11 -17 23c-9-4-18-6-28-7l-4-28h-21l-4 28c-10 1-19 4-28 7l-17-23 -18 11 11 26c-8 6-14 13-20 20l-26-11 -11 18 23 17c-4 9-6 18-7 28l-28 4v21l28 4c1 10 4 19 7 28l-23 17 11 18 26-11c6 8 13 14 20 20l-11 26 18 11 17-23c9 4 18 6 28 7l4 28h21l4-28c10-1 19-4 28-7l17 23 18-11 -11-26c8-6 14-13 20-20l26 11 11-18 -23-17c4-9 6-18 7-28L392 390zM255 453c-41 0-74-33-74-74 0-41 33-74 74-74 41 0 74 33 74 74C328 420 295 453 255 453z'/>       </g>     </svg> </div>"

                scoreDiv.innerHTML += loading_gears;
                //awaitResults(scoreDiv, path);
            }

            if (scores.message.includes('cached')) {
                scoreDiv.innerHTML += '<br><br>'
                scoreDiv.innerHTML += '<button class=requestScore id=requestButton> Request Score </button>'

                document.getElementById('requestButton').addEventListener('click', async function() {
                    scoreDiv.innerHTML = "<h2 class=\"h4 mb-3\"> OSS Scores </h2>"

                    //maybe put this in try catch?
                    requestScoreCalculation(owner, repo, 'activity');
                    requestScoreCalculation(owner, repo, 'license');
                    scoreDiv.innerHTML += 'Scores Requested';
                    awaitResults(scoreDiv, owner, repo);
                });
            }
        }
    })

}

async function awaitResults(scoreDiv, owner, repo) {
    promiseTimeout(2500).then(() => {
        console.log('Requesting Score');
        getScores(owner, repo).then(scores => {
            message = scores.message;
            if (message == undefined|| message == '') {
                insertScores(scoreDiv, scores);
            } else {
                awaitResults(scoreDiv, owner, repo);
            }
        });
    });
}

async function getFakeScores(owner, repo) {
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

async function getScores(owner, repo) {
    let scores = {license: null, activity: null, message: null};
    let promises = [];    
    let licenseRequestUrl = basePath + '/owner/' + owner + '/repo/' + repo + '/type/license';
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

    let activityRequestUrl = basePath + '/owner/' + owner + '/repo/' + repo + '/type/activity';
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

    insertScoreSection(owner, repo, scoreDiv, getScores(owner, repo, scoreDiv));
    // insertScoreSection(owner, repo, scoreDiv, getFakeScores(owner, repo, scoreDiv));
}
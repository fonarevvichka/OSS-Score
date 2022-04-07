// Imports
import React, {useState} from 'react'
import './Homepage.css';
import DisplayScores from './DisplayScores.js';
import { AiOutlineInfoCircle } from "react-icons/ai";

/* functional component for homepage */
export default function Home(props) {

    const [inputs, setInputs] = useState("");
    //const basePath = 'https://hvacjx4u1l.execute-api.us-east-2.amazonaws.com/prod/catalog/github' //prod
    const basePath = 'https://xvzhkajkzh.execute-api.us-east-2.amazonaws.com/dev/catalog/github' //dev
    // const basePath = 'https://4oam7avy4i.execute-api.us-east-2.amazonaws.com/staging/catalog/github' //staging
    const calculationMessages = ['Metric not yet calculated', 'Error querying score', 'Data out of date']
    //const errorMessages = ['Could not access repo, check that it was inputted correctly and is public', 'Cannot provide score for private repo']
    const loading_gears = "<div class='loading-extension'>     <svg class='machine-extension'xmlns='http://www.w3.org/2000/svg' x='0px' y='0px' viewBox='0 0 645 526'>       <defs/>       <g>         <path  x='-173,694' y='-173,694' class='large-shadow-extension' d='M645 194v-21l-29-4c-1-10-3-19-6-28l25-14 -8-19 -28 7c-5-8-10-16-16-24L602 68l-15-15 -23 17c-7-6-15-11-24-16l7-28 -19-8 -14 25c-9-3-18-5-28-6L482 10h-21l-4 29c-10 1-19 3-28 6l-14-25 -19 8 7 28c-8 5-16 10-24 16l-23-17L341 68l17 23c-6 7-11 15-16 24l-28-7 -8 19 25 14c-3 9-5 18-6 28l-29 4v21l29 4c1 10 3 19 6 28l-25 14 8 19 28-7c5 8 10 16 16 24l-17 23 15 15 23-17c7 6 15 11 24 16l-7 28 19 8 14-25c9 3 18 5 28 6l4 29h21l4-29c10-1 19-3 28-6l14 25 19-8 -7-28c8-5 16-10 24-16l23 17 15-15 -17-23c6-7 11-15 16-24l28 7 8-19 -25-14c3-9 5-18 6-28L645 194zM471 294c-61 0-110-49-110-110S411 74 471 74s110 49 110 110S532 294 471 294z'/>       </g>       <g>         <path x='-136,996' y='-136,996' class='medium-shadow-extension' d='M402 400v-21l-28-4c-1-10-4-19-7-28l23-17 -11-18L352 323c-6-8-13-14-20-20l11-26 -18-11 -17 23c-9-4-18-6-28-7l-4-28h-21l-4 28c-10 1-19 4-28 7l-17-23 -18 11 11 26c-8 6-14 13-20 20l-26-11 -11 18 23 17c-4 9-6 18-7 28l-28 4v21l28 4c1 10 4 19 7 28l-23 17 11 18 26-11c6 8 13 14 20 20l-11 26 18 11 17-23c9 4 18 6 28 7l4 28h21l4-28c10-1 19-4 28-7l17 23 18-11 -11-26c8-6 14-13 20-20l26 11 11-18 -23-17c4-9 6-18 7-28L402 400zM265 463c-41 0-74-33-74-74 0-41 33-74 74-74 41 0 74 33 74 74C338 430 305 463 265 463z'/>       </g>       <g >         <path x='-100,136' y='-100,136' class='small-shadow-extension' d='M210 246v-21l-29-4c-2-10-6-18-11-26l18-23 -15-15 -23 18c-8-5-17-9-26-11l-4-29H100l-4 29c-10 2-18 6-26 11l-23-18 -15 15 18 23c-5 8-9 17-11 26L10 225v21l29 4c2 10 6 18 11 26l-18 23 15 15 23-18c8 5 17 9 26 11l4 29h21l4-29c10-2 18-6 26-11l23 18 15-15 -18-23c5-8 9-17 11-26L210 246zM110 272c-20 0-37-17-37-37s17-37 37-37c20 0 37 17 37 37S131 272 110 272z'/>       </g>       <g>         <path x='-100,136' y='-100,136' class='small-extension' d='M200 236v-21l-29-4c-2-10-6-18-11-26l18-23 -15-15 -23 18c-8-5-17-9-26-11l-4-29H90l-4 29c-10 2-18 6-26 11l-23-18 -15 15 18 23c-5 8-9 17-11 26L0 215v21l29 4c2 10 6 18 11 26l-18 23 15 15 23-18c8 5 17 9 26 11l4 29h21l4-29c10-2 18-6 26-11l23 18 15-15 -18-23c5-8 9-17 11-26L200 236zM100 262c-20 0-37-17-37-37s17-37 37-37c20 0 37 17 37 37S121 262 100 262z'/>       </g>       <g>         <path x='-173,694' y='-173,694' class='large-extension' d='M635 184v-21l-29-4c-1-10-3-19-6-28l25-14 -8-19 -28 7c-5-8-10-16-16-24L592 58l-15-15 -23 17c-7-6-15-11-24-16l7-28 -19-8 -14 25c-9-3-18-5-28-6L472 0h-21l-4 29c-10 1-19 3-28 6L405 9l-19 8 7 28c-8 5-16 10-24 16l-23-17L331 58l17 23c-6 7-11 15-16 24l-28-7 -8 19 25 14c-3 9-5 18-6 28l-29 4v21l29 4c1 10 3 19 6 28l-25 14 8 19 28-7c5 8 10 16 16 24l-17 23 15 15 23-17c7 6 15 11 24 16l-7 28 19 8 14-25c9 3 18 5 28 6l4 29h21l4-29c10-1 19-3 28-6l14 25 19-8 -7-28c8-5 16-10 24-16l23 17 15-15 -17-23c6-7 11-15 16-24l28 7 8-19 -25-14c3-9 5-18 6-28L635 184zM461 284c-61 0-110-49-110-110S401 64 461 64s110 49 110 110S522 284 461 284z'/>       </g>       <g>         <path x='-136,996' y='-136,996' class='medium-extension' d='M392 390v-21l-28-4c-1-10-4-19-7-28l23-17 -11-18L342 313c-6-8-13-14-20-20l11-26 -18-11 -17 23c-9-4-18-6-28-7l-4-28h-21l-4 28c-10 1-19 4-28 7l-17-23 -18 11 11 26c-8 6-14 13-20 20l-26-11 -11 18 23 17c-4 9-6 18-7 28l-28 4v21l28 4c1 10 4 19 7 28l-23 17 11 18 26-11c6 8 13 14 20 20l-11 26 18 11 17-23c9 4 18 6 28 7l4 28h21l4-28c10-1 19-4 28-7l17 23 18-11 -11-26c8-6 14-13 20-20l26 11 11-18 -23-17c4-9 6-18 7-28L392 390zM255 453c-41 0-74-33-74-74 0-41 33-74 74-74 41 0 74 33 74 74C328 420 295 453 255 453z'/>       </g>     </svg> </div>";
    let scoreDisplay = ''

    /* function to validate github URL. Returns true if valid, false otherwise */
    const validateURL = (url, repoNum) => {

        // Pattern match Github URL with regex: check that it starts with https://github.co
        // Note: If there is a user error in the input of the repository name, which would be a valid
        //      repository name (ex. github.comfacebook/react) our validation will pass
        //      but the API will return no repo found

        const validgitHub = new RegExp('^(https://)?github.com/+[a-zA-Z0-9._-]+/+[a-zA-Z0-9._-]+$')
        const validgitHubTree = new RegExp('^(https://)?github.com/+[a-zA-Z0-9._-]+/+[a-zA-Z0-9._-]+/tree/+[a-zA-Z0-9._-]+$')

        const validgitHubNoPrefix = new RegExp('^(?!github.com/)[a-zA-Z0-9._-]+/+[a-zA-Z0-9._-]+$')
        const validgitHubTreeNoPrefix = new RegExp('^(?!github.com/)[a-zA-Z0-9._-]+/+[a-zA-Z0-9._-]+/tree/+[a-zA-Z0-9._-]+$')
        
        if (!validgitHub.test(url) && !validgitHubTree.test(url) &&
            !validgitHubNoPrefix.test(url) && !validgitHubTreeNoPrefix.test(url)) {
            return false
        }

        hideError(repoNum)
        return true;
    }

    /* function to display error for invalid github URL */
    const displayError = (repoNum) => {
        // changing css
        document.getElementById("search" + repoNum).style.borderColor = "#cc0000"  
        document.getElementById("error-message" + repoNum).style.visibility = "visible"
    }

    /* function to hide error */
    const hideError = (repoNum) => {
        document.getElementById("search" + repoNum).style.borderColor = "#000000"
        document.getElementById("error-message" + repoNum).style.visibility = "hidden"
    }


    /* handleChange function for inputs */
    const handleChange = (repoNum) => (event) => {
        
        // updating target value
        if (repoNum === "1") {
            setInputs(state => ({ ...state, search1: event.target.value }))
        } else if (repoNum === "2") {
            setInputs(state => ({ ...state, search2: event.target.value }))
        }

        // displaying errors if invalid
        if (!validateURL(event.target.value, repoNum) && event.target.value !== "") {
            displayError(repoNum)
        } else {
            hideError(repoNum)
        }
    }


    /* function for parsing name and author */
    const getNameAuthor = (url) => {
        // or replace github.com/
        let newUrl = null;

        if (url.includes('https://github.com/')) {
            newUrl = url.replace('https://github.com/', '');
        } else if (url.includes('github.com/')) {
            newUrl = url.replace('github.com/', '')
        } else {
            newUrl = url
        }

        let splitUrl = newUrl.split('/').filter(element => element !== '');
        let owner = '';
        let repo = '';

        if (splitUrl.length === 2) { // Repo homepage
            owner = splitUrl[0];
            repo = splitUrl[1];
        } else if (splitUrl.length === 4) {
            if (splitUrl[2] === 'tree') { // on specific branch
                owner = splitUrl[0];
                repo = splitUrl[1];
            }
        }
        return [owner, repo];
    }

    async function promiseTimeout(time) {
        return new Promise(function (resolve, reject) {
            setTimeout(function () { resolve(); }, time);
        });
    };

    async function requestScores(owner, repo) {
        let message = null;
        let success = false
        let requestURL = basePath + '/owner/' + owner + '/name/' + repo;
        let promise =
            fetch(requestURL, {
                method: 'POST',
                mode: 'cors'
            }).then(async (response) => {
                if (response.status === 201) {
                    let messagePromise = response.json();
                    await messagePromise.then(response => {
                        message = response.message;
                        success = true
                    }).catch(err => {
                        console.error(err);
                    });
                } else if (response.status === 406) {
                    message = "Cannot provide score for private repo";
                } else if ((response.status === 501) || (response.status === 503)) {
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



    /* function that makes api call given an owner and repo name, returns metrics in json */
    const getMetrics = async (owner, repo) => {
        let metric_name = 'all'
        try {
            let response = await fetch(basePath + '/owner/' + owner + '/name/' + repo + '/metric/'
                    + metric_name)
            if (response.status === 200) {
                let data = await response.json()
                if (data.message === "Metric ready") {
                    return [owner, repo, data];
                } else if (calculationMessages.includes(data.message)) {
                    let requestResponse = await requestScores(owner, repo)
                    if (requestResponse.success) {
                        let metrics = await awaitResults(owner, repo)
                        return [owner, repo, metrics];
                    }
                } else {
                    // score calculate queued, don't call request scores, go straight to awaitResults
                    let metrics = await awaitResults(owner, repo)
                    return [owner, repo, metrics];
                }
            } else if (response.status === 406) {
                alert(owner + '/' + repo + " is private or does not exist")
                console.error("Repository entered does not exist")
            } else {
                alert("Error connecting to OSS-Score API")
                console.error("Error connecting to OSS-Score API")
            }
        } catch (error) {
            console.error(error)
            return Promise.reject(error);
        }
    }


    async function awaitResults(owner, repo) {
        let metric_name = 'all'
        let response = await fetch(basePath + '/owner/' + owner + '/name/' + repo + '/metric/'
                + metric_name)
        if (response.status === 200) {
            let data = await response.json()
            if (data.message === "Metric ready") {
                return data
            } else {
                await promiseTimeout(1000)
                return awaitResults(owner, repo)
            }
        }
        return null
    }

    /* handleSubmit function that does everything */
    const handleSubmit = async (evt) => {
        
        // Disable button to prevent multiple clicks
        document.getElementById("compare-button").disabled = true;

        // erorr prevent default
        evt.preventDefault();

        // Clear all html in head2head and hide it
        document.getElementById("head2head").innerHTML = ''
        scoreDisplay = ''
        
        // Show loading gear
        document.getElementById("loading").innerHTML += loading_gears;

        let scorePromises = []

        if (validateURL(inputs.search1, "1")) {
            // parse Name and Author, call API
            let [owner1, name1] = getNameAuthor(inputs.search1)
            scorePromises.push(getMetrics(owner1, name1))
        } else {
            displayError("1");
        }

        if (validateURL(inputs.search2, "2")) {
            // parse Name and Author, call API
            let [owner2, name2] = getNameAuthor(inputs.search2)
            scorePromises.push(getMetrics(owner2, name2))
        } else {
            displayError("2");
        }

        await Promise.all(scorePromises).then((values) => {
            scoreDisplay += DisplayScores(values)
        }).catch(e => console.log('Error caught', e));

 
        // Hide loading gear/clear all html in head2head
        document.getElementById("loading").innerHTML = ''
        document.getElementById("head2head").innerHTML = ''

        document.getElementById("head2head").innerHTML += scoreDisplay

        // Enable button
        document.getElementById("compare-button").disabled = false;
    }

    return (
        <div className="Home">
            <div className="logo"></div>
            <form onSubmit={handleSubmit}>
                <div class="searchbar">
                    <div>
                        <label htmlFor="search1" >Link to Github repo #1</label><br></br>
                        <input key="search1" id="search1" name="search1" type="text" placeholder="Search Repo 1" onClick={() => document.getElementById('search1').style.borderColor = '#000000'}
                            onChange={handleChange("1")} value={inputs.search1}/>
                        <div class="tool-tip-repo"> <AiOutlineInfoCircle color="white" />
                            <span class="tooltiptext-repo" style={{ width: "230px", marginLeft: "-115px" }}> <div>Insert github repo as:</div>
                                <div>owner/name</div><div>github.com/owner/name</div><div>https://github.com/owner/name</div>
                            </span>
                        </div>
                        <div class="error-message" id="error-message1" name="error-message1">Please enter a valid Github URL</div>
                    </div>
                    <div>
                        <label htmlFor="search2" >Link to Github repo #2</label><br></br>
                        <input key="search2" id="search2" name="search2" type="text" placeholder="Search Repo 2" onClick={() => document.getElementById('search2').style.borderColor = '#000000'}
                            onChange = {handleChange("2")} value={inputs.search2} />
                        <div class="tool-tip-repo"> <AiOutlineInfoCircle color="white" />
                             <span class="tooltiptext-repo" style={{ width: "230px", marginLeft: "-115px" }}> <div>Insert github repo as:</div>
                                <div>owner/name</div><div>github.com/owner/name</div><div>https://github.com/owner/name</div>
                            </span>
                        </div>
                        <div class="error-message" id="error-message2" name="error-message2">Please enter a valid Github URL</div>
                    </div>
                </div>
                <div class="compare">
                    <button id="compare-button" class="compare-button" type="submit" value="Submit">Get Metrics</button>
                </div>
            </form>
            <div id="loading"></div>
            <div class="head2head" id="head2head"></div>

            < svg view-box="0 0 1600 900" >
                <path fill="#5b4693" opacity="1" d="M0,297C267,403,534,368,801,577,C1068,786,1335,391,1602,451,C1600, 900,1600, 900,1600, 900C1600, 900,1600, 900,1600, 900C1600, 900,1600, 900,1600, 900C1600, 900,1600, 900,1600, 900L1600,900C1333,900,1066,900,799,900,C532,900,265,900,-2,900,C0, 900,0, 900,0, 900C0, 900,0, 900,0, 900C0, 900,0, 900,0, 900C0, 900,0, 900,0, 900L1401,900L0,900Z" />
            </svg >
        </div>
    );
}
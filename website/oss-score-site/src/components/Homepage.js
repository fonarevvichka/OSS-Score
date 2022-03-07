
// imports
import { data } from 'jquery';
import React, {useState} from 'react'
import './Homepage.css';
import DisplayScores from './DisplayScores.js';




/* functional component for homepage */
export default function Home(props) {

    const [inputs, setInputs] = useState("");
    const basePath = 'https://ru8ibij7yc.execute-api.us-east-2.amazonaws.com/staging/catalog/github'
    const loading_gears = "<div class='loading-extension'>     <svg class='machine-extension'xmlns='http://www.w3.org/2000/svg' x='0px' y='0px' viewBox='0 0 645 526'>       <defs/>       <g>         <path  x='-173,694' y='-173,694' class='large-shadow-extension' d='M645 194v-21l-29-4c-1-10-3-19-6-28l25-14 -8-19 -28 7c-5-8-10-16-16-24L602 68l-15-15 -23 17c-7-6-15-11-24-16l7-28 -19-8 -14 25c-9-3-18-5-28-6L482 10h-21l-4 29c-10 1-19 3-28 6l-14-25 -19 8 7 28c-8 5-16 10-24 16l-23-17L341 68l17 23c-6 7-11 15-16 24l-28-7 -8 19 25 14c-3 9-5 18-6 28l-29 4v21l29 4c1 10 3 19 6 28l-25 14 8 19 28-7c5 8 10 16 16 24l-17 23 15 15 23-17c7 6 15 11 24 16l-7 28 19 8 14-25c9 3 18 5 28 6l4 29h21l4-29c10-1 19-3 28-6l14 25 19-8 -7-28c8-5 16-10 24-16l23 17 15-15 -17-23c6-7 11-15 16-24l28 7 8-19 -25-14c3-9 5-18 6-28L645 194zM471 294c-61 0-110-49-110-110S411 74 471 74s110 49 110 110S532 294 471 294z'/>       </g>       <g>         <path x='-136,996' y='-136,996' class='medium-shadow-extension' d='M402 400v-21l-28-4c-1-10-4-19-7-28l23-17 -11-18L352 323c-6-8-13-14-20-20l11-26 -18-11 -17 23c-9-4-18-6-28-7l-4-28h-21l-4 28c-10 1-19 4-28 7l-17-23 -18 11 11 26c-8 6-14 13-20 20l-26-11 -11 18 23 17c-4 9-6 18-7 28l-28 4v21l28 4c1 10 4 19 7 28l-23 17 11 18 26-11c6 8 13 14 20 20l-11 26 18 11 17-23c9 4 18 6 28 7l4 28h21l4-28c10-1 19-4 28-7l17 23 18-11 -11-26c8-6 14-13 20-20l26 11 11-18 -23-17c4-9 6-18 7-28L402 400zM265 463c-41 0-74-33-74-74 0-41 33-74 74-74 41 0 74 33 74 74C338 430 305 463 265 463z'/>       </g>       <g >         <path x='-100,136' y='-100,136' class='small-shadow-extension' d='M210 246v-21l-29-4c-2-10-6-18-11-26l18-23 -15-15 -23 18c-8-5-17-9-26-11l-4-29H100l-4 29c-10 2-18 6-26 11l-23-18 -15 15 18 23c-5 8-9 17-11 26L10 225v21l29 4c2 10 6 18 11 26l-18 23 15 15 23-18c8 5 17 9 26 11l4 29h21l4-29c10-2 18-6 26-11l23 18 15-15 -18-23c5-8 9-17 11-26L210 246zM110 272c-20 0-37-17-37-37s17-37 37-37c20 0 37 17 37 37S131 272 110 272z'/>       </g>       <g>         <path x='-100,136' y='-100,136' class='small-extension' d='M200 236v-21l-29-4c-2-10-6-18-11-26l18-23 -15-15 -23 18c-8-5-17-9-26-11l-4-29H90l-4 29c-10 2-18 6-26 11l-23-18 -15 15 18 23c-5 8-9 17-11 26L0 215v21l29 4c2 10 6 18 11 26l-18 23 15 15 23-18c8 5 17 9 26 11l4 29h21l4-29c10-2 18-6 26-11l23 18 15-15 -18-23c5-8 9-17 11-26L200 236zM100 262c-20 0-37-17-37-37s17-37 37-37c20 0 37 17 37 37S121 262 100 262z'/>       </g>       <g>         <path x='-173,694' y='-173,694' class='large-extension' d='M635 184v-21l-29-4c-1-10-3-19-6-28l25-14 -8-19 -28 7c-5-8-10-16-16-24L592 58l-15-15 -23 17c-7-6-15-11-24-16l7-28 -19-8 -14 25c-9-3-18-5-28-6L472 0h-21l-4 29c-10 1-19 3-28 6L405 9l-19 8 7 28c-8 5-16 10-24 16l-23-17L331 58l17 23c-6 7-11 15-16 24l-28-7 -8 19 25 14c-3 9-5 18-6 28l-29 4v21l29 4c1 10 3 19 6 28l-25 14 8 19 28-7c5 8 10 16 16 24l-17 23 15 15 23-17c7 6 15 11 24 16l-7 28 19 8 14-25c9 3 18 5 28 6l4 29h21l4-29c10-1 19-3 28-6l14 25 19-8 -7-28c8-5 16-10 24-16l23 17 15-15 -17-23c6-7 11-15 16-24l28 7 8-19 -25-14c3-9 5-18 6-28L635 184zM461 284c-61 0-110-49-110-110S401 64 461 64s110 49 110 110S522 284 461 284z'/>       </g>       <g>         <path x='-136,996' y='-136,996' class='medium-extension' d='M392 390v-21l-28-4c-1-10-4-19-7-28l23-17 -11-18L342 313c-6-8-13-14-20-20l11-26 -18-11 -17 23c-9-4-18-6-28-7l-4-28h-21l-4 28c-10 1-19 4-28 7l-17-23 -18 11 11 26c-8 6-14 13-20 20l-26-11 -11 18 23 17c-4 9-6 18-7 28l-28 4v21l28 4c1 10 4 19 7 28l-23 17 11 18 26-11c6 8 13 14 20 20l-11 26 18 11 17-23c9 4 18 6 28 7l4 28h21l4-28c10-1 19-4 28-7l17 23 18-11 -11-26c8-6 14-13 20-20l26 11 11-18 -23-17c4-9 6-18 7-28L392 390zM255 453c-41 0-74-33-74-74 0-41 33-74 74-74 41 0 74 33 74 74C328 420 295 453 255 453z'/>       </g>     </svg> </div>";


    /* function to validate github URL. Returns true if valid, false otherwise */
    const validateURL = (url, repoNum) => {

        // Pattern match Github URL with regex: check that it starts with https://github.com
        const validgitHub = new RegExp('^https://github.com/+[a-zA-Z0-9._-]+/+[a-zA-Z0-9._-]+$')
        const validgitHubTree = new RegExp('^https://github.com/+[a-zA-Z0-9._-]+/+[a-zA-Z0-9._-]+/tree/+[a-zA-Z0-9._-]+$')
        
        if (!validgitHub.test(url) && !validgitHubTree.test(url)) {
            return false
        }

        document.getElementById("search" + repoNum).style.borderColor = "#000000"
        document.getElementById("error-message" + repoNum).style.display = "none"
        return true;
    }

    /* function to display error for invalid github URL */
    const displayError = (repoNum) => {
        // changing css
        document.getElementById("search" + repoNum).style.borderColor = "#cc0000"  
        document.getElementById("error-message" + repoNum).style.display = "block"
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
        if (!validateURL(event.target.value, repoNum)) {
            displayError(repoNum)
        }
    }

    /*function handleFocusOut(event) {
        alert("in focus out")
        if (!validateURL(event.target.value, "1")) {
            displayError("1")
        }
    }*/


    /* function for parsing name and author */
    const getNameAuthor = (url) => {
        let newUrl = url.replace('https://github.com/', '');

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


    function promiseTimeout(time) {
        return new Promise(function (resolve, reject) {
            setTimeout(function () { resolve(); }, time);
        });
    };

    async function requestScores(owner, repo) {
        let message = null;
        let requestURL = basePath + '/owner/' + owner + '/name/' + repo;
        let promise =
            fetch(requestURL, {
                method: 'POST',
                mode: 'cors'
            }).then(async (response) => {
                if (response.status == 200) {
                    let messagePromise = response.json();
                    await messagePromise.then(response => {
                        message = response.message;
                    }).catch(err => {
                        console.error(err);
                    });
                } else {
                    let messagePromise = response.json();
                    await messagePromise.then(response => {
                        message = response.message;
                    }).catch(err => {
                        console.error(err);
                    });
                }
            }).catch(err => {
                message = "error caught in post";
                message += err;
                console.error(err);
            });
        await promise;
        return message;
    }


    /* function that makes api call given an owner and repo name, returns metrics in json */
    const getMetrics = async (owner_name, repo_name) => {
        let catalog_name = 'github'
        let metric_name = 'all'
        try {
            let response = await fetch('https://ru8ibij7yc.execute-api.us-east-2.amazonaws.com/staging/catalog/'
                    + catalog_name + '/owner/' + owner_name + '/name/' + repo_name + '/metric/'
                    + metric_name)
            
            if (response.status == 200) {
                return response.json()
            } else if (response.status == 406) {
                console.error("Repository entered does not exist")
                return null
            } else {
                console.error("Error connecting to OSS-Score API")
                return null
            }
        } catch (error) {
            console.log(error)
            return [];
        }
    }

    async function awaitResults(owner_name, repo_name) {
        promiseTimeout(500).then(() => {
            console.log('Requesting Score');
            getMetrics(owner_name, repo_name).then(scores => {
                if (scores.activity != null) {
                    // need to call get metrics maybe
                    // insertScores(scores); // Call component to make html
                    console.log("Got the scores finally")
                } else {
                    // Loading gears
                    awaitResults(owner_name, repo_name);
                }
            });
        });
    }



    /* function for flattenning JSON */
    /*const flattenJSON = (obj = {}, res = {}) => {
        // base case: when it is not an object
        // base case: when it is an object with a score and a confidence

        // recursive case: when it is an object that doesn't have score and confidence

        for (var key in obj) {
            if (typeof obj[key] !== 'object') {
                res[key] = obj[key];
            } else if (obj[key].hasOwnProperty('message') && obj[key].hasOwnProperty('metric') && obj[key].hasOwnProperty('confidence')) {
                res[key] = obj[key];
            } else {
                flattenJSON(obj[key], res);
            };
        };
        
        return res;
    }*/

    /* handleSubmit function that does everything */
    const handleSubmit = async (evt) => {

        // erorr prevent default
        evt.preventDefault();

        // Hide loading gear/clear all html in head2head
        document.getElementById("head2head").innerHTML = ''

        // validating first url
        let owner1 = null
        let name1 = null
        let scores1 = null
        
        let owner2 = null
        let name2 = null
        let scores2 = null

        // Loading gear
        document.getElementById("head2head").innerHTML += loading_gears;

        /*if (!validgitHub.test(inputs.search1) && !validgitHubTree.test(inputs.search1)) {
            // clear textbox and highlight textbox red
            document.getElementById("search1").style.borderColor = "#cc0000"
        } else {
            // parse Name and Author, call API
            [owner1, name1] = getNameAuthor(inputs.search1)
            scores1 = await getMetrics(owner1, name1)

        }*/

        if (validateURL(inputs.search1, "1")) {
            // parse Name and Author, call API
            [owner1, name1] = getNameAuthor(inputs.search1)
            scores1 = await getMetrics(owner1, name1)
        } else {
            displayError("1");
        }

        // validating second url 
        /*if (!validgitHub.test(inputs.search2) && !validgitHubTree.test(inputs.search2)) {
            // clear textbox and highlight textbox red
            document.getElementById("search2").style.borderColor = "#cc0000"
        } else {
            // parse Name and Author, call API
            [owner2, name2] = getNameAuthor(inputs.search2)
            scores2 = await getMetrics(owner2, name2)

            // Hide loading gear
            //document.getElementById("loading-extension").style.display = 'none'
        }*/

        if (validateURL(inputs.search2, "2")) {
            // parse Name and Author, call API
            [owner2, name2] = getNameAuthor(inputs.search2)
            scores2 = await getMetrics(owner2, name2)
        } else {
            displayError("2");
        }

        // Hide loading gear/clear all html in head2head
        document.getElementById("head2head").innerHTML = ''
        document.getElementById("head2head").style.display = 'flex'
        if (scores1 != null && scores2 != null) {
            // Display both scores
            document.getElementById("head2head").innerHTML += DisplayScores(owner1, name1, scores1)
            document.getElementById("head2head").innerHTML += DisplayScores(owner2, name2, scores2)
        } else if (scores1 != null) {
            // Display score 1
            document.getElementById("head2head").innerHTML += DisplayScores(owner1, name1, scores1)
        } else if (scores2 != null) {
            // Display score 2
            document.getElementById("head2head").innerHTML += DisplayScores(owner2, name2, scores2)
        }

        //document.getElementById("loading").innerHTML = loading_gears;


        // get name and author: https://github.com/author/name
        // const [owner1, name1] = getNameAuthor(inputs.search1)
        // const [owner2, name2] = getNameAuthor(inputs.search2)

        // set name and author in html
        // document.getElementById("repoName1").innerHTML = name1;
        // document.getElementById("repoAuthor1").innerHTML = owner1;
        // document.getElementById("repoName2").innerHTML = name2;
        // document.getElementById("repoAuthor2").innerHTML = owner2;
        


        // make call to api to get metrics in teh form of json
        // let scores1 = await getMetrics(owner1, name1)

        // if (scores1.activity != null) { // VALID SCORES RETURNED
        //     //insertScores(scores); // our component that writes the html
        // } else {
        //     console.log("Requesting previously unknown score");
        //     //submit post request
        //     awaitResults(owner1, name1)
        // }


        // let scores2 = await getMetrics(owner2, name2)
        
        //console.log(scores1)
        //console.log(scores2)


        // // call component to display metrics
        // // iterate through metric div tags to put scores in
        // var num_repos = 2;
        // var metricElems = document.getElementsByClassName("metric");
        // var confidenceElems = document.getElementsByClassName("confidence")

        // for (var i = 0; i < metricElems.length; i++) {
        //     var metricId = metricElems[i].id
        //     metricId = metricId.slice(0, -1)
        //     console.log(metricId)

        //     if (i < (metricElems.length / num_repos)) {
        //         metricElems[i].innerHTML = scores1[metricId].metric
        //         confidenceElems[i].innerHTML = scores1[metricId].confidence
        //     }
        //     else {
        //         metricElems[i].innerHTML = scores2[metricId].metric
        //         confidenceElems[i].innerHTML = scores2[metricId].confidence
        //     }
        // }

        // show scores
        // document.getElementById("loading").style.display = 'none'
        // document.getElementById("head2head").style.display = 'flex'

        // // clear previous scores
        // document.getElementById("head2head").innerHTML = ""
        // document.getElementById("head2head").innerHTML = ""

        // // display new scores
        // document.getElementById("head2head").innerHTML += DisplayScores(owner1, name1, scores1)
        // document.getElementById("head2head").innerHTML += DisplayScores(owner2, name2, scores2)
    }

    return (
        <div className="Home">
            {/*<img src="../images/logo1.png" alt="OSS-SCORE"></img>*/}
            <header>OSS-SCORE</header>
            <form onSubmit={handleSubmit}>
                <div class="searchbar">
                    <div>
                        <label for="search1" >Link to Github repo #1</label><br></br>
                        {/*<input key="search1" id="search1" name="search1" type="text" placeholder="Search Repo 1" onClick={() => document.getElementById('search1').style.borderColor = '#000000'}
                            onChange={({ target }) => setInputs(state => ({ ...state, search1: target.value }))} value={inputs.search1} />*/}
                        <input key="search1" id="search1" name="search1" type="text" placeholder="Search Repo 1" onClick={() => document.getElementById('search1').style.borderColor = '#000000'}
                            onChange={handleChange("1")} value={inputs.search1}/>
                        <div class="error-message" id="error-message1" name="error-message1">Please enter a valid Github URL</div>
                    </div>
                    <div>
                        <label for="search2" >Link to Github repo #2</label><br></br>
                        {/* <input key="search2" id="search2" name="search2" type="text" placeholder="Search Repo 2" onClick={() => document.getElementById('search2').style.borderColor = '#000000'}
                            onChange={({ target }) => setInputs(state => ({ ...state, search2: target.value }))} value={inputs.search2} /> */}
                        <input key="search2" id="search2" name="search2" type="text" placeholder="Search Repo 2" onClick={() => document.getElementById('search2').style.borderColor = '#000000'}
                            onChange = {handleChange("2")} value={inputs.search2} />
                        <div class="error-message" id="error-message2" name="error-message2">Please enter a valid Github URL</div>
                    </div>
                </div>
                <div class="compare">
                    <button class="compare-button" type="submit" value="Submit">Compare</button>
                </div>
            </form>
            <div id="loading"></div>
            <div class="head2head" id="head2head"></div>
        </div>
    );
}
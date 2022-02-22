
// imports
import { data } from 'jquery';
import React, {useState} from 'react'
import './Homepage.css';


/* functional component for homepage */
export default function Home(props) {

    const [inputs, setInputs] = useState("");

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

    /* function that makes api call given an owner and repo name, returns metrics in json */
    const getMetrics = async (owner_name, repo_name) => {
        let catalog_name = 'github'
        let metric_name = 'all'
        try {
            let response = await fetch('https://ru8ibij7yc.execute-api.us-east-2.amazonaws.com/staging/catalog/'
                    + catalog_name + '/owner/' + owner_name + '/name/' + repo_name + '/metric/'
                    + metric_name)
            
            return response.json()
        } catch (error) {
            console.log(error)
            return [];
        }
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
        evt.preventDefault();

        /* css for displaying html on submit (not working) */
        if (document.getElementById("head2head").style.display === 'none') {
            document.getElementById("head2head").style.display = 'block'
        }

        // Pattern match Github URL with regex: check that it starts with https://github.com
        const validgitHub = new RegExp('^https://github.com/+[a-zA-Z0-9._-]+/+[a-zA-Z0-9._-]+$')
        const validgitHubTree = new RegExp('^https://github.com/+[a-zA-Z0-9._-]+/+[a-zA-Z0-9._-]+/tree/+[a-zA-Z0-9._-]+$')

        // validating first url
        let repo1isValid = true 

        if (!validgitHub.test(inputs.search1) && !validgitHubTree.test(inputs.search1)) {
            // clear textbox and highlight textbox red
            document.getElementById("search1").style.borderColor = "#cc0000"
            document.getElementById("search1").value = ''
            repo1isValid = false
        }

        // validating second url 
        let repo2isValid = true
        if (!validgitHub.test(inputs.search2) && !validgitHubTree.test(inputs.search2)) {
            // clear textbox and highlight textbox red
            document.getElementById("search2").style.borderColor = "#cc0000"
            document.getElementById("search2").value = ''
            repo2isValid = false
        }

        // exit if either is invalid
        if (!(repo1isValid && repo2isValid)) {
            return
        }

        // get name and author: https://github.com/author/name
        const [owner1, name1] = getNameAuthor(inputs.search1)
        const [owner2, name2] = getNameAuthor(inputs.search2)

        // set name and author in html
        document.getElementById("repoName1").innerHTML = name1;
        document.getElementById("repoAuthor1").innerHTML = owner1;
        document.getElementById("repoName2").innerHTML = name2;
        document.getElementById("repoAuthor2").innerHTML = owner2;
        


        // make call to api to get metrics in teh form of json
        let scores1 = await getMetrics(owner1, name1)
        let scores2 = await getMetrics(owner2, name2)
        
        //console.log(scores1)
        //console.log(scores2)


        // call component to display metrics
        // iterate through metric div tags to put scores in
        var num_repos = 2;
        var metricElems = document.getElementsByClassName("metric");
        var confidenceElems = document.getElementsByClassName("confidence")

        for (var i = 0; i < metricElems.length; i++) {
            var metricId = metricElems[i].id
            metricId = metricId.slice(0, -1)
            console.log(metricId)

            if (i < (metricElems.length / num_repos)) {
                metricElems[i].innerHTML = scores1[metricId].metric
                confidenceElems[i].innerHTML = scores1[metricId].confidence
            }
            else {
                metricElems[i].innerHTML = scores2[metricId].metric
                confidenceElems[i].innerHTML = scores2[metricId].confidence
            }
        }



        //document.getElementById("score-first").innerHTML = inputs.search1 + inputs.search2;
    }

    return (
        <div className="Home">
            <img src="../images/logo1.png" alt="OSS-SCORE"></img>
            <header>OSS-SCORE</header>
            <form onSubmit={handleSubmit}>
                <div class="searchbar">
                    <div>
                        <label for="search1" >Link to Github repo #1</label><br></br>
                        <input key="search1" id="search1" name="search1" type="text" placeholder="Search Repo 1" onClick={() => document.getElementById('search1').style.borderColor = '#000000'}
                            onChange={({ target }) => setInputs(state => ({ ...state, search1: target.value }))} value={inputs.search1} />
                    </div>
                    <div>
                        <label for="search2" >Link to Github repo #2</label><br></br>
                        <input key="search2" id="search2" name="search2" type="text" placeholder="Search Repo 2" onClick={() => document.getElementById('search2').style.borderColor = '#000000'}
                            onChange={({ target }) => setInputs(state => ({ ...state, search2: target.value }))} value={inputs.search2} />
                    </div>
                </div>
                <div class="compare">
                    <button class="compare-button" type="submit" value="Submit">Compare</button>
                </div>
            </form>
            <div class="head2head" id="head2head">
                <div class="repo1-stats">
                    <div class="basic-info-display">
                        <div class="basic-info-title">Name</div>
                        <div class="basic-info" id="repoName1">None</div>
                        <div class="basic-info-title">Author</div>
                        <div class="basic-info" id="repoAuthor1">None</div>
                    </div>
                    {/* <div class="repo-header">
                        <div class="metric-container-title">Overall Score</div>
                        <div class="metric" id="overallScore1">0</div>
                        <div class="confidence" id="overallConfScore1">Confidence: 0</div>
                    </div> */}
                    <div class="metrics-display">
                        <div class="metric-category">Activity Scores</div>
                        <div class="metric-container">
                            <div class="metric-container-title">Activity Score</div>
                            <div class="metric" id="repoActivityScore1">0</div>
                            <div class="confidence" id="repoActivityConfScore1">Confidence: 0</div>
                        </div>
                        <div class="metric-container">
                            <div class="metric-container-title">Issue Closure Time</div>
                            <div class="metric" id="issueClosureTime1">0</div>
                            <div class="confidence" id="issueClosureTimeConf1">Confidence: 0</div>
                        </div>
                        
                        <div class="metric-container">
                            <div class="metric-container-title">Commit Cadence</div>
                            <div class="metric" id="commitCadence1">0</div>
                            <div class="confidence" id="commitCadenceConf1">Confidence: 0</div>
                        </div>


                        <div class="metric-container">
                            <div class="metric-container-title">Release Score</div>
                            {/* <div class="metric" id="releaseScore1">0</div>
                            <div class="confidence" id="releaseConfScore1">Confidence: 0</div> */}
                            <div class="submetric-container">
                                <div class="submetric-container-title">Age Last Release</div>
                                <div class="metric" id="ageLastRelease1">0</div>
                                <div class="confidence" id="ageLastReleaseConf1">Confidence: 0</div>
                            </div>
                            <div class="submetric-container">
                                <div class="submetric-container-title">Release Cadence</div>
                                <div class="metric" id="releaseCadence1">0</div>
                                <div class="confidence" id="releaseCadenceConf1">Confidence: 0</div>
                            </div>
                        </div>

                        <div class="metric-container">
                            <div class="metric-container-title">Contributors</div>
                            <div class="metric" id="contributors1">0</div>
                            <div class="confidence" id="contributorsConf1">Confidence: 0</div>
                        </div>
                    </div>
                    <div class="repo-licence-score">
                        <div class="metric-category">License Scores</div>
                        <div class="metric-container">
                            <div class="metric-container-title">License Score</div>
                            <div class="metric" id="repoLicenseScore1">0</div>
                            <div class="confidence" id="repoLicenseConfScore1">Confidence: 0</div>
                        </div>
                    </div>
                    <div class="repo-dependency-score">
                        <div class="metric-category">Dependency Scores</div>
                        <div class="metric-container">
                            <div class="metric-container-title">Dependency Activity Score</div>
                            <div class="metric" id="dependencyActivityScore1">0</div>
                            <div class="confidence" id="dependencyActivityConfScore1">Confidence: 0</div>
                        </div>
                        <div class="metric-container">
                            <div class="metric-container-title">Dependency License Score</div>
                            <div class="metric" id="dependencyLicenseScore1">0</div>
                            <div class="confidence" id="dependencyLicenseConfScore1">Confidence: 0</div>
                        </div>
                    </div>
                    <div class="repo-stars">
                        <div class="metric-category">Stars</div>
                        <div class="metric-container">
                            <div class="metric-container-title">Stars</div>
                            <div class="metric" id="stars1">0</div>
                            <div class="confidence" id="stars1">Confidence: 0</div>
                        </div>
                    </div>
                </div>
                <div class="repo2-stats">
                    <div class="basic-info-display">
                        <div class="basic-info-title">Name</div>
                        <div class="basic-info" id="repoName2">None</div>
                        <div class="basic-info-title">Author</div>
                        <div class="basic-info" id="repoAuthor2">None</div>
                    </div>
                    {/* <div class="repo-header">
                        <div class="metric-container-title">Overall Score</div>
                        <div class="metric" id="overallScore1">0</div>
                        <div class="confidence" id="overallConfScore1">Confidence: 0</div>
                    </div> */}
                    <div class="metrics-display">
                        <div class="metric-category">Activity Scores</div>
                        <div class="metric-container">
                            <div class="metric-container-title">Activity Score</div>
                            <div class="metric" id="repoActivityScore2">0</div>
                            <div class="confidence" id="repoActivityConfScore2">Confidence: 0</div>
                        </div>
                        <div class="metric-container">
                            <div class="metric-container-title">Issue Closure Time</div>
                            <div class="metric" id="issueClosureTime2">0</div>
                            <div class="confidence" id="issueClosureTimeConf2">Confidence: 0</div>
                        </div>
                        
                        <div class="metric-container">
                            <div class="metric-container-title">Commit Cadence</div>
                            <div class="metric" id="commitCadence2">0</div>
                            <div class="confidence" id="commitCadenceConf2">Confidence: 0</div>
                        </div>


                        <div class="metric-container">
                            <div class="metric-container-title">Release Score</div>
                            {/* <div class="metric" id="releaseScore2">0</div>
                            <div class="confidence" id="releaseConfScore2">Confidence: 0</div> */}
                            <div class="submetric-container">
                                <div class="submetric-container-title">Age Last Release</div>
                                <div class="metric" id="ageLastRelease2">0</div>
                                <div class="confidence" id="ageLastReleaseConf2">Confidence: 0</div>
                            </div>
                            <div class="submetric-container">
                                <div class="submetric-container-title">Release Cadence</div>
                                <div class="metric" id="releaseCadence2">0</div>
                                <div class="confidence" id="releaseCadenceConf2">Confidence: 0</div>
                            </div>
                        </div>

                        <div class="metric-container">
                            <div class="metric-container-title">Contributors</div>
                            <div class="metric" id="contributors2">0</div>
                            <div class="confidence" id="contributorsConf2">Confidence: 0</div>
                        </div>
                    </div>
                    <div class="repo-licence-score">
                        <div class="metric-category">License Scores</div>
                        <div class="metric-container">
                            <div class="metric-container-title">License Score</div>
                            <div class="metric" id="repoLicenseScore2">0</div>
                            <div class="confidence" id="repoLicenseConfScore2">Confidence: 0</div>
                        </div>
                    </div>
                    <div class="repo-dependency-score">
                        <div class="metric-category">Dependency Scores</div>
                        <div class="metric-container">
                            <div class="metric-container-title">Dependency Activity Score</div>
                            <div class="metric" id="dependencyActivityScore2">0</div>
                            <div class="confidence" id="dependencyActivityConfScore2">Confidence: 0</div>
                        </div>
                        <div class="metric-container">
                            <div class="metric-container-title">Dependency License Score</div>
                            <div class="metric" id="dependencyLicenseScore2">0</div>
                            <div class="confidence" id="dependencyLicenseConfScore2">Confidence: 0</div>
                        </div>
                    </div>
                    <div class="repo-stars">
                        <div class="metric-category">Stars</div>
                        <div class="metric-container">
                            <div class="metric-container-title">Stars</div>
                            <div class="metric" id="stars2">0</div>
                            <div class="confidence" id="stars2">Confidence: 0</div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
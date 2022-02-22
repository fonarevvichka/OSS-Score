import { data } from 'jquery';
import React, {useState} from 'react'
//import { useForm } from 'react-hook-form';

import './Homepage.css';

/* functional component for homepage */
export default function Home(props) {
    // const {register, handleSubmit} = useForm();

    const [inputs, setInputs] = useState("");

    
    // const onSubmit = (data) => {
    //     console.log(data)
    // }

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

    /* function for flattenning JSON */
    const flattenJSON = (obj = {}, res = {}) => {
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
    }


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


    // const componentDidMount = (owner_name, repo_name) => {

    //     let catalog_name = 'github'
    //     let metric_name = 'all'
    //     // GET request using fetch with error handling
    //     fetch('https://ru8ibij7yc.execute-api.us-east-2.amazonaws.com/staging/catalog/' + catalog_name + '/owner/' + owner_name + '/name/' + repo_name + '/metric/' + metric_name)
    //         .then(async response => {
    //             const data = await response.json();

    //             // check for error response
    //             if (!response.ok) {
    //                 // get error message from body or default to response statusText
    //                 const error = (data && data.message) || response.statusText;
    //                 return Promise.reject(error);
    //             }
 
    //             // this.setState({ totalReactPackages: data.total })
    //         })
    //         .catch(error => {
    //             this.setState({ errorMessage: error.toStrin() });
    //             console.error('There was an error!', error);
    //         });
    //     return data;
    // }

    // const testing_git = () => {
        // note: based on playing around with repo name, only non-alphanumeric characters is -, _, and .
        // assuming author has same restrictions

        // let ghtests = ["http://github.com/elidow/oss-score", // fails 0
        //             "https://githubb.com/elidow/oss-score", // 1
        //             "https://github.co/elidow/oss-score", // 2
        //             "https://github.com/elidow", // 3
        //             "https://github.com/elidow****/oss", // 4
        //             "https://github.com/&elidow/oss", // 5
        //             "https://github.com/elidow/oss***", // 6 
        //             "https://github.com/elidow/***oss", // 7
        //             "https://github.com/elidow/oss/yes", // 8 
        //             "https://github.com/elidow/oss/score/", // 9 
        //             "https://github.com/elidow///oss", // 10 
        //             "https://github.com/elidow/oss",  // 11     // succeeds
        //             "https://github.com/elidow/oss-score", // 12
        //             "https://github.com/elidow123/oss", // 13 
        //             "https://github.com/eli_dow./oss-score" // 14
        //             ]

        // //let answers = [false] * 11 + [true] * 4;
        // let answers = [false, false, false, false, false, false, false, false, false, false, false, true, true, true, true]
        // for (let i = 0; i < tests.length; i++) {
        //     if ((validgitHub.test(tests[i]) && !answers[i]) || ((!validgitHub.test(tests[i]) && answers[i]))) {
        //         alert(`Test: ${i} failed`)
        //     }
        // }
    //}


    const handleSubmit = async (evt) => {
        evt.preventDefault();

        if (document.getElementById("head2head").style.display === 'none') {
            document.getElementById("head2head").style.display = 'block'
        }

        // First do website validation

        // Pattern match URL
        // First validate its a githubx url: by checking that it starts with https://github.com
        const validgitHub = new RegExp('^https://github.com/+[a-zA-Z0-9._-]+/+[a-zA-Z0-9._-]+$')
        const validgitHubTree = new RegExp('^https://github.com/+[a-zA-Z0-9._-]+/+[a-zA-Z0-9._-]+/tree/+[a-zA-Z0-9._-]+$')

        let repo1isValid = true 
        if (!validgitHub.test(inputs.search1) && !validgitHubTree.test(inputs.search1)) {
            // clear textbox and highlight textbox red
            document.getElementById("search1").style.borderColor = "#cc0000"
            document.getElementById("search1").value = ''
            repo1isValid = false
        }

        let repo2isValid = true
        if (!validgitHub.test(inputs.search2) && !validgitHubTree.test(inputs.search2)) {
            // clear textbox and highlight textbox red
            document.getElementById("search2").style.borderColor = "#cc0000"
            document.getElementById("search2").value = ''
            repo2isValid = false
        }

        if (!(repo1isValid && repo2isValid)) {
            return
        }

        // Next, extract author, https://github.com/author/name
        //call getnameauthor

        const [owner1, name1] = getNameAuthor(inputs.search1)
        const [owner2, name2] = getNameAuthor(inputs.search2)

        // set author and repo names in html
        document.getElementById("repoName1").innerHTML = name1;
        document.getElementById("repoAuthor1").innerHTML = owner1;
        document.getElementById("repoName2").innerHTML = name2;
        document.getElementById("repoAuthor2").innerHTML = owner2;
        


        // Here: call the api to get metrics for inputs.search1 and inputs.search2


        let scores1 = await getMetrics(owner1, name1)
        let scores2 = await getMetrics(owner2, name2)
        
        console.log(scores1)
        console.log(scores2)


        // Sample JSON object

        // JSON object names must be the same as the id's in the html
        // let scores1 = {
        //     "overallScore": {"score": 90, "confidence": 80},
        //     "activityScore": {
        //         "activityScore": {"score": 75, "confidence": 50},
        //         "commitScore": {"score": 75, "confidence": 12},
        //         "contributorScore": {"score": 85, "confidence": 62},
        //         "releaseScore": {
        //             "releaseScore":{"score": 20, "confidence": 1},
        //             "ageLastReleaseScore":{"score":3, "confidence": 100},
        //             "releaseCadenceScore": { "score": 40, "confidence": 10 }
        //             },
        //         "issueScore": {"score":99, "confidence": 99}
        //     },
        //     "licenseScore": {"score": 3, "confidence": 100},
        //     "dependencyActivityScore": {"score": 100, "confidence": 2},
        //     "dependencyLicenseScore": {"score": 3, "confidence": 100},
        // }

        // console.log(flattenJSON(scores));

        //alert("sup");

        var num_repos = 2;

        // let scores1 = {
        //     "overallScore": { "score": 90, "confidence": 80 },
        //     "activityScore": { "score": 75, "confidence": 50 },
        //     "commitScore": { "score": 75, "confidence": 12 },
        //     "contributorScore": { "score": 85, "confidence": 62 },
        //     "releaseScore": { "score": 20, "confidence": 1 },
        //     "ageLastReleaseScore": { "score": 3, "confidence": 100 },
        //     "releaseCadenceScore": { "score": 40, "confidence": 10 },
        //     "issueScore": { "score": 99, "confidence": 99 },
        //     "licenseScore": { "score": 3, "confidence": 100 },
        //     "dependencyActivityScore": { "score": 100, "confidence": 2 },
        //     "dependencyLicenseScore": { "score": 3, "confidence": 100 }
        // }

        // let scores2 = {
        //     "overallScore": { "score": 100, "confidence": 80 },
        //     "activityScore": { "score": 75, "confidence":90 },
        //     "commitScore": { "score": 99, "confidence": 82 },
        //     "contributorScore": { "score": 8, "confidence": 100 },
        //     "releaseScore": { "score": 20, "confidence": 100 },
        //     "ageLastReleaseScore": { "score": 70, "confidence": 100 },
        //     "releaseCadenceScore": { "score": 80, "confidence": 10 },
        //     "issueScore": { "score": 45, "confidence": 99 },
        //     "licenseScore": { "score": 3, "confidence": 100 },
        //     "dependencyActivityScore": { "score": 100, "confidence": 2 },
        //     "dependencyLicenseScore": { "score": 3, "confidence": 100 }
        // }


        // scores1 = flattenJSON(scores1)
        // console.log(scores1)
        // scores2 = flattenJSON(scores2)
        // console.log(scores2)

        // Once we have metrics: call component to display metrics

        // Iterate through to put scores in
        var metricElems = document.getElementsByClassName("metric");
        var confidenceElems = document.getElementsByClassName("confidence")

        for (var i = 0; i < metricElems.length; i++) {
            var metricId = metricElems[i].id
            metricId = metricId.slice(0, -1)
            console.log(metricId)

            //alert(scoreId)
            if (i < (scoreElems.length / num_repos)) {
                metricElems[i].innerHTML = scores1[metricId].metric
                confidenceElems[i].innerHTML = scores1[metricId].confidence
            }
            else {
                scoreElems[i].innerHTML = scores2[scoreId].metric
                confidenceElems[i].innerHTML = scores2[scoreId].confidence
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
                <div class="repo1">
                    <div class="repo-header">
                        <div class="stat-subHeader">Name</div>
                        <div class="repo-name" id="repoName1">None</div>
                        <div class="stat-subHeader">Author</div>
                        <div class="repo-name" id="repoAuthor1">None</div>
                    </div>
                    {/* <div class="repo-header">
                        <div class="stat-Header">Overall Score</div>
                        <div class="metric" id="overallScore1">0</div>
                        <div class="confidence" id="overallConfScore1">Confidence: 0</div>
                    </div> */}

                    <div class="repo-header">
                        <div class="subheaderTitle">Activity Scores</div>
                        <div class="repo-subheader">
                            <div class="stat-Header">Activity Score</div>
                            <div class="metric" id="repoActivityScore1">0</div>
                            <div class="confidence" id="repoActivityConfScore1">Confidence: 0</div>
                        </div>


                        <div class="repo-subheader">
                            <div class="stat-Header">Issue Closure Time</div>
                            <div class="metric" id="issueClosureTime1">0</div>
                            <div class="confidence" id="issueClosureTimeConf1">Confidence: 0</div>
                        </div>
                        
                        <div class="repo-subheader">
                            <div class="stat-Header">Commit Cadence</div>
                            <div class="metric" id="commitCadence1">0</div>
                            <div class="confidence" id="commitCadenceConf1">Confidence: 0</div>
                        </div>


                        <div class="repo-subheader">
                            <div class="stat-Header">Release Score</div>
                            {/* <div class="metric" id="releaseScore1">0</div>
                            <div class="confidence" id="releaseConfScore1">Confidence: 0</div> */}
                            <div class="stat-subheader">
                                <div class="scoreLabel">Age Last Release</div>
                                <div class="metric" id="ageLastRelease1">0</div>
                                <div class="confidence" id="ageLastReleaseConf1">Confidence: 0</div>
                                <div class="scoreLabel">Release Cadence</div>
                                <div class="metric" id="releaseCadence1">0</div>
                                <div class="confidence" id="releaseCadenceConf1">Confidence: 0</div>
                            </div>
                        </div>

                        <div class="repo-subheader">
                            <div class="stat-Header">Contributors</div>
                            <div class="metric" id="contributors1">0</div>
                            <div class="confidence" id="contributorsConf1">Confidence: 0</div>
                        </div>
                    </div>
                    <div class="repo-licence-score">
                        <div class="subheaderTitle">License Scores</div>
                        <div class="repo-subheader">
                            <div class="stat-Header">License Score</div>
                            <div class="metric" id="repoLicenseScore1">0</div>
                            <div class="confidence" id="repoLicenseConfScore1">Confidence: 0</div>
                        </div>
                    </div>
                    <div class="repo-dependency-score">
                        <div class="subheaderTitle">Dependency Scores</div>
                        <div class="repo-subheader">
                            <div class="stat-Header">Dependency Activity Score</div>
                            <div class="metric" id="dependencyActivityScore1">0</div>
                            <div class="confidence" id="dependencyActivityConfScore1">Confidence: 0</div>
                        </div>
                        <div class="repo-subheader">
                            <div class="stat-Header">Dependency License Score</div>
                            <div class="metric" id="dependencyLicenseScore1">0</div>
                            <div class="confidence" id="dependencyLicenseConfScore1">Confidence: 0</div>
                        </div>
                    </div>
                    <div class="repo-stars">
                        <div class="subheaderTitle">Stars</div>
                        <div class="repo-subheader">
                            <div class="stat-Header">Stars</div>
                            <div class="metric" id="stars1">0</div>
                            <div class="confidence" id="stars1">Confidence: 0</div>
                        </div>
                    </div>
                </div>

                <div class="repo2">
                <div class="repo-header">
                        <div class="stat-subHeader">Name</div>
                        <div class="repo-name" id="repoName2">None</div>
                        <div class="stat-subHeader">Author</div>
                        <div class="repo-name" id="repoAuthor2">None</div>
                    </div>
                    <div class="repo-header">
                        <div class="stat-Header">Overall Score</div>
                        <div class="metric" id="overallScore2">0</div>
                        <div class="confidence" id="overallConfScore2">Confidence: 0</div>
                    </div>
                    
                    <div class="repo-header">
                        <div class="subheaderTitle">Activity Scores</div>
                        <div class="repo-subheader">
                            <div class="stat-Header">Activity Score</div>
                            <div class="metric" id="activityScore2">0</div>
                            <div class="confidence" id="activityConfScore2">Confidence: 0</div>
                        </div>
                        
                        <div class="repo-subheader">
                            <div class="stat-Header">Commit Score</div>
                            <div class="metric" id="commitScore2">0</div>
                            <div class="confidence" id="commitConfScore2">Confidence: 0</div>
                        </div>
                        
                        <div class="repo-subheader">
                            <div class="stat-Header">Contributor Score</div>
                            <div class="metric" id="contributorScore2">0</div>
                            <div class="confidence" id="contributorConfScore2">Confidence: 0</div>
                        </div>
                        
                        <div class="repo-subheader">
                            <div class="stat-Header">Release Score</div>
                            <div class="metric" id="releaseScore2">0</div>
                            <div class="confidence" id="releaseConfScore2">Confidence: 0</div>
                            <div class="stat-subheader">
                            <div class="scoreLabel">Age Last Release Score</div>
                                <div class="metric" id="ageLastReleaseScore2">0</div>
                                <div class="confidence" id="ageLastReleaseConfScore2">Confidence: 0</div>
                                <div class="scoreLabel">Release Cadence Score</div>
                                <div class="metric" id="releaseCadenceScore2">0</div>
                                <div class="confidence" id="releaseCadenceConfScore2">Confidence: 0</div>
                            </div>
                        </div>
                        
                        <div class="repo-subheader">
                            <div class="stat-Header">Issue Closure Time Score</div>
                            <div class="metric" id="issueScore2">0</div>
                            <div class="confidence" id="issueConfScore2">Confidence: 0</div>
                        </div>
                    </div>
                    <div class="repo-licence-score">
                        <div class="subheaderTitle">License Scores</div>
                        <div class="repo-subheader">
                            <div class="stat-Header">License Score</div>
                            <div class="metric" id="licenseScore2">0</div>
                            <div class="confidence" id="licenseConfScore2">Confidence: 0</div>
                        </div>
                    </div>
                    <div class="repo-dependency-score">
                        <div class="subheaderTitle">Dependency Scores</div>
                        <div class="repo-subheader">
                            <div class="stat-Header">Dependency Activity Score</div>
                            <div class="metric" id="dependencyActivityScore2">0</div>
                            <div class="confidence" id="dependencyActivityConfScore2">Confidence: 0</div>
                        </div>
                        <div class="repo-subheader">
                            <div class="stat-Header">Dependency License Score</div>
                            <div class="metric" id="dependencyLicenseScore2">0</div>
                            <div class="confidence" id="dependencyLicenseConfScore2">Confidence: 0</div>
                        </div>
                    </div>
                    <div class="repo-stars">
                        <div class="subheaderTitle">Stars</div>
                        <div class="repo-subheader">
                            <div class="stat-Header">Stars</div>
                            <div class="metric" id="stars2">0</div>
                            <div class="confidence" id="stars2">Confidence: 0</div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
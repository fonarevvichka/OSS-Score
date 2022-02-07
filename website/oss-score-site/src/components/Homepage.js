import React, {useState} from 'react'
//import { useForm } from 'react-hook-form';

import './Homepage.css';


export default function Home(props) {
    // const {register, handleSubmit} = useForm();

    const [inputs, setInputs] = useState("");

    
    // const onSubmit = (data) => {
    //     console.log(data)
    // }

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


    const handleSubmit = (evt) => {
        evt.preventDefault();

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
        
        //alert(`Submitting Repos: ${inputs.search1} and ${inputs.search2}`)

        // Sample JSON object

        // JSON object names must be the same as the id's in the html
        // let scores = {
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
        //     "dependencyLicenseScore": {"score": 3, "confidence": 100}
        // }

        alert("sup");

        let scores = {
            "overallScore": { "score": 90, "confidence": 80 },
            "activityScore": { "score": 75, "confidence": 50 },
            "commitScore": { "score": 75, "confidence": 12 },
            "contributorScore": { "score": 85, "confidence": 62 },
            "releaseScore": { "score": 20, "confidence": 1 },
            "ageLastReleaseScore": { "score": 3, "confidence": 100 },
            "releaseCadenceScore": { "score": 40, "confidence": 10 },
            "issueScore": { "score": 99, "confidence": 99 },
            "licenseScore": { "score": 3, "confidence": 100 },
            "dependencyActivityScore": { "score": 100, "confidence": 2 },
            "dependencyLicenseScore": { "score": 3, "confidence": 100 }
        }

        let scores2 = {
            "overallScore": { "score": 100, "confidence": 80 },
            "activityScore": { "score": 75, "confidence":90 },
            "commitScore": { "score": 99, "confidence": 82 },
            "contributorScore": { "score": 8, "confidence": 100 },
            "releaseScore": { "score": 20, "confidence": 100 },
            "ageLastReleaseScore": { "score": 70, "confidence": 100 },
            "releaseCadenceScore": { "score": 80, "confidence": 10 },
            "issueScore": { "score": 45, "confidence": 99 },
            "licenseScore": { "score": 3, "confidence": 100 },
            "dependencyActivityScore": { "score": 100, "confidence": 2 },
            "dependencyLicenseScore": { "score": 3, "confidence": 100 }
        }
        // Once we have metrics: call component to display metrics

        // Iterate through to put scores in
        var scoreElems = document.getElementsByClassName("score");
        var confidenceElems = document.getElementsByClassName("confidence")

        for (var i = 0; i < scoreElems.length; i++) {
            var scoreId = scoreElems[i].id
            scoreId = scoreId.slice(0, -1)

            //alert(scoreId)
            if (i < (scoreElems.length / 2)) {
                scoreElems[i].innerHTML = scores[scoreId].score
                confidenceElems[i].innerHTML = scores[scoreId].confidence
            }
            else {
                scoreElems[i].innerHTML = scores2[scoreId].score
                confidenceElems[i].innerHTML = scores2[scoreId].confidence
            }
        }





        //document.getElementById("score-first").innerHTML = inputs.search1 + inputs.search2;
    }

    return (
        <div className="Home">
            <img src="../../../images/logo1.png" alt="OSS-SCORE"></img>
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
            <div class="head2head">
                <div class="repo1">
                    <div class="repo-header">
                        <div class="stat-subHeader">Name</div>
                        <div class="repo-name" id="repoName1">None</div>
                        <div class="stat-subHeader">Author</div>
                        <div class="repo-name" id="repoAuthor1">None</div>
                    </div>
                    <div class="repo-header">
                        <div class="stat-Header">Overall Score</div>
                        <div class="score" id="overallScore1">0</div>
                        <div class="confidence" id="overallConfScore1">Confidence: 0</div>
                    </div>

                    <div class="repo-header">
                        <div class="subheaderTitle">Activity Score</div>
                        <div class="repo-subheader">
                            <div class="stat-Header">Activity Score</div>
                            <div class="score" id="activityScore1">0</div>
                            <div class="confidence" id="activityConfScore1">Confidence: 0</div>
                        </div>

                        <div class="repo-subheader">
                            <div class="stat-Header">Commit Score</div>
                            <div class="score" id="commitScore1">0</div>
                            <div class="confidence" id="commitConfScore1">Confidence: 0</div>
                        </div>

                        <div class="repo-subheader">
                            <div class="stat-Header">contributor Score</div>
                            <div class="score" id="contributorScore1">0</div>
                            <div class="confidence" id="contributorConfScore1">Confidence: 0</div>
                        </div>

                        <div class="repo-subheader">
                            <div class="stat-Header">Release Score</div>
                            <div class="score" id="releaseScore1">0</div>
                            <div class="confidence" id="releaseConfScore1">Confidence: 0</div>
                            <div class="stat-subheader">
                                <div class="scoreLabel">Age Last Release Score</div>
                                <div class="score" id="ageLastReleaseScore1">0</div>
                                <div class="confidence" id="ageLastReleaseConfScore1">Confidence: 0</div>
                                <div class="scoreLabel">Release Cadence Score</div>
                                <div class="score" id="releaseCadenceScore1">0</div>
                                <div class="confidence" id="releaseCadenceConfScore1">Confidence: 0</div>
                            </div>
                        </div>

                        <div class="repo-subheader">
                            <div class="stat-Header">Issue Closure Time Score</div>
                            <div class="score" id="issueScore1">0</div>
                            <div class="confidence" id="issueConfScore1">Confidence: 0</div>
                        </div>
                    </div>
                    <div class="repo-licence-score">
                        <div class="subheaderTitle">License Score</div>
                        <div class="repo-subheader">
                            <div class="stat-Header">License Score</div>
                            <div class="score" id="licenseScore1">0</div>
                            <div class="confidence" id="licenseConfScore1">Confidence: 0</div>
                        </div>
                    </div>
                    <div class="repo-dependency-score">
                        <div class="subheaderTitle">Dependency Score</div>
                        <div class="repo-subheader">
                            <div class="stat-Header">Dependency Activity Score</div>
                            <div class="score" id="dependencyActivityScore1">0</div>
                            <div class="confidence" id="dependencyActivityConfScore1">Confidence: 0</div>
                        </div>
                        <div class="repo-subheader">
                            <div class="stat-Header">Dependency License Score</div>
                            <div class="score" id="dependencyLicenseScore1">0</div>
                            <div class="confidence" id="dependencyLicenseConfScore1">Confidence: 0</div>
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
                        <div class="score" id="overallScore2">0</div>
                        <div class="confidence" id="overallConfScore2">Confidence: 0</div>
                    </div>
                    
                    <div class="repo-header">
                        <div class="subheaderTitle">Activity Score</div>
                        <div class="repo-subheader">
                            <div class="stat-Header">Activity Score</div>
                            <div class="score" id="activityScore2">0</div>
                            <div class="confidence" id="activityConfScore2">Confidence: 0</div>
                        </div>
                        
                        <div class="repo-subheader">
                            <div class="stat-Header">Commit Score</div>
                            <div class="score" id="commitScore2">0</div>
                            <div class="confidence" id="commitConfScore2">Confidence: 0</div>
                        </div>
                        
                        <div class="repo-subheader">
                            <div class="stat-Header">contributor Score</div>
                            <div class="score" id="contributorScore2">0</div>
                            <div class="confidence" id="contributorConfScore2">Confidence: 0</div>
                        </div>
                        
                        <div class="repo-subheader">
                            <div class="stat-Header">Release Score</div>
                            <div class="score" id="releaseScore2">0</div>
                            <div class="confidence" id="releaseConfScore2">Confidence: 0</div>
                            <div class="stat-subheader">
                            <div class="scoreLabel">Age Last Release Score</div>
                                <div class="score" id="ageLastReleaseScore2">0</div>
                                <div class="confidence" id="ageLastReleaseConfScore2">Confidence: 0</div>
                                <div class="scoreLabel">Release Cadence Score</div>
                                <div class="score" id="releaseCadenceScore2">0</div>
                                <div class="confidence" id="releaseCadenceConfScore2">Confidence: 0</div>
                            </div>
                        </div>
                        
                        <div class="repo-subheader">
                            <div class="stat-Header">Issue Closure Time Score</div>
                            <div class="score" id="issueScore2">0</div>
                            <div class="confidence" id="issueConfScore2">Confidence: 0</div>
                        </div>
                    </div>
                    <div class="repo-licence-score">
                        <div class="subheaderTitle">License Score</div>
                        <div class="repo-subheader">
                            <div class="stat-Header">License Score</div>
                            <div class="score" id="licenseScore2">0</div>
                            <div class="confidence" id="licenseConfScore2">Confidence: 0</div>
                        </div>
                    </div>
                    <div class="repo-dependency-score">
                        <div class="subheaderTitle">Dependency Score</div>
                        <div class="repo-subheader">
                            <div class="stat-Header">Dependency Activity Score</div>
                            <div class="score" id="dependencyActivityScore2">0</div>
                            <div class="confidence" id="dependencyActivityConfScore2">Confidence: 0</div>
                        </div>
                        <div class="repo-subheader">
                            <div class="stat-Header">Dependency License Score</div>
                            <div class="score" id="dependencyLicenseScore2">0</div>
                            <div class="confidence" id="dependencyLicenseConfScore2">Confidence: 0</div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
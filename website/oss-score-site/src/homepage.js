import React from 'react'

import './homepage.css';


export default function Home(){
    return (       
        <div className="Home">
            <header>OSS-SCORE</header>
            <div class="searchbar">
                <div>
                    <label for="search1" >Link to Github repo #1</label><br></br>
                    <input name="search1" type="text" placeholder="Search Repo 1"></input>
                </div>
                <div>
                    <label for="search2" >Link to Github repo #2</label><br></br>
                    <input name="search2" type="text" placeholder="Search Repo 2"></input>
                </div>
            </div>
            <div class="compare">
                <button class="compare-button">Compare</button>
            </div>
            
            <div className="head2head">
                <div name="repo1">
                    <div class="repo-header">
                        <div class="repo-name"></div>
                        <div class="repo-activity-score">
                            <div class="score"></div>
                            <div class="confidence"></div>
                        </div>

                        <div class="repo-lisence-score">
                            <div class="score"></div>
                            <div class="confidence"></div>
                        </div>
                    </div>

                    <div class="issue-closure">
                        <div class="stat-header"></div>
                        <div class="stat"></div>
                    </div>
                    <div class="age-last-release"></div>
                </div>

                <div name="repo2">
                    <div class="repo-header">
                        <div class="repo-name"></div>
                        <div class="repo-activity-score">
                            <div class="score"></div>
                            <div class="confidence"></div>
                        </div>

                        <div class="repo-lisence-score">
                            <div class="score"></div>
                            <div class="confidence"></div>
                        </div>
                    </div>

                    <div class="issue-closure">
                        <div class="stat-header"></div>
                        <div class="stat"></div>
                    </div>
                    <div class="age-last-release">
                        <div class="stat-header"></div>
                        <div class="stat"></div>
                    </div>

                </div>
                
            </div>
            
        </div>
    );
}
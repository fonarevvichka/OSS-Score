import React from 'react'
import { FaLinkedin } from "react-icons/fa";
import './About.css';


const About = () => {
  return (
    <div className='about-container'>
      <header className='main-header'>About</header>
      <div className='opener'>Have you ever had trouble finding the right open source software package for your project? Have you ever started using a package only to realize hours later that it has long been deprecated or does not have active developer groups to fix any outstanding issues?</div>
      <header>Mission</header>
      <div className='mission-statement'>Our mission at OSS-Score is to help developers make quick and informed decisions about OSS tools. We want to minimize the time it takes from deciding you need to use an open source package to integrating it into your project. In addition to creating a more streamlined workflow experience, we also want to make sure developers are fully aware of all the details about a package before they decide to use it.</div>
      <header>Our Solution</header>
      <div className='solution'>
        <div>We provide developers with a comprehensive activity and license score, that is calculated from a wide variety of GitHub metadata for that repository and its dependencies. This ensures developers can quickly and easily pick repositories that are well maintained and themselves rely on well maintained dependencies.</div>
        <div>For quick "at-a-glance" information the <a className='internal-link' href='/extension'>OSS-Score chrome extension</a> embeds scores directly into the GitHub repository homepage.</div>
        <div class="parent">
          <img class="image1 " alt='extension-img' src={require('../images/extension_pic.png')} />
          <img class="image2" alt='extension-closeup-img' src={require('../images/extension_closeup_pic.png')} />
        </div>
        <div>For a deeper dive the <a className='internal-link' href='/'>OSS-Score website</a> allows users to compare the scores and individual metrics of different open source projects and tools.</div>
        <div className='image-container'><img alt='website-img' src={require('../images/website_pic.png')}></img></div>

        <p className="sectionHeader">
          Disclaimers
        </p>
        <p className="disclaimers">
          - All metrics need to be pulled from GitHub, we try to cache as much as possible but sometimes queries can take a long time if we have no data cached.
          New repos can take anywhere from 30 seconds to a few minutes. Please be patient.
          <br/>
          <br/>
          - We have an aritificial shelf life for our data. Data is considiered 'in-date' if is is less than three days old.
          As such some metrics may not quite line up with what you see on the repo homepage.
        </p>

        <p className="sectionHeader">
          Feedback
        </p>
      
        <p className="Feedback">
        We welcome any and all feedback! Have some ideas about the scoring algorithm, new metrics we should track? Please join us on our GitHub page.
        File some issues, start a discussion, we can't wait to hear from you.
        </p>

        <p className="sectionHeader">
          Next Steps (or maybe just pondering)
        </p>
        <p className="nextSteps">
          - Adding authentication to our lambda endpoints. We are thinking of using something like OAUTH tokens combined with API Gateways internal services to autheticate the requests coming in. <br/>
          - We want to add pull requests closure rate as an input to our scoring calculation<br/>
          - Incorporate the recency of certain data so that they have a greater impact on the score <br/>
          - Exploring some GitHub security metrics as an input <br/>
        </p>
      </div>

      <header>Meet the Team</header>
      <div className='meet-the-team'>
        <div className='team-blurb'>OSS-Score was created by Vichka Fonarev, Eli Dow, Daniel Alderman, and Emil Polakiewicz as a senior computer science capstone project at Tufts University. Open source enthusiast and current CTO of Kubeshop Ole Lensmar provided invaluable mentorship throughout the creation of OSS-Score.</div>  
        <br />
        <div className='team-members'>
            <div className='team-member'>
              <div className='name'>Vichka Fonarev</div>
            <div className='photo'><img alt='vichka-pic' src={require('../images/vichka.png')}></img></div>
            <div className='linkedin'><a target="_blank" rel="noreferrer" href='https://www.linkedin.com/in/vichka-fonarev-b1b980110/'><FaLinkedin /></a></div>
            </div>

            <div className='team-member'>
              <div className='name'>Eli Dow</div>
            <div className='photo'><img alt='eliIsABum' src={require('../images/eli.png')}></img></div>
            <div className='linkedin'><a target="_blank" rel="noreferrer" href='https://www.linkedin.com/in/eli-dow-93105b168/'><FaLinkedin /></a></div>
            </div>

            <div className='team-member'>
              <div className='name'>Daniel Alderman</div>
            <div className='photo'><img alt='doniel-pic' src={require('../images/daniel.png')}></img></div>
            <div className='linkedin'><a target="_blank" rel="noreferrer" href='https://www.linkedin.com/in/daniel-alderman-ab88321a1/'><FaLinkedin /></a></div>
            </div>

            <div className='team-member'>
              <div className='name'>Emil Polakiewicz</div>
            <div className='photo'><img alt='emil-pic' src={require('../images/emil.png')}></img></div>
            <div className='linkedin'><a target="_blank" rel="noreferrer" href='https://www.linkedin.com/in/emil-polakiewicz-12887b17b/'><FaLinkedin /></a></div>
            </div>
        </div>
      </div>
    </div>
  );
}

export default About;
import React from 'react'
import { FaLinkedin } from "react-icons/fa";
import './About.css';


const About = () => {
  return (
    <div className='about-container'>
      <header className='main-header'>About</header>
      <p className='opener'>&emsp; Have you ever had trouble finding the right open source software package for your project? Have you ever started using a package only to realize hours later that it has long been deprecated or does not have active developer groups to fix any outstanding issues?</p>
      <header>Mission</header>
      <p className='mission-statement'>&emsp; Our mission at OSS-Score is to help developers make quick and informed decisions about OSS tools. We want to minimize the time it takes from deciding you need to use an open source package to integrating it into your project. In addition to creating a more streamlined workflow experience, we also want to make sure developers are fully aware of all the details about a package before they decide to use it.</p>
      <header>Our Solution</header>
      <p className='solution'>
        <p>&emsp; We provide developers with a comprehensive activity and license score, that is calculated from a wide variety of GitHub metadata for that repository and its dependencies. This ensures developers can quickly and easily pick repositories that are well maintained and won’t cause any headaches in the future.</p>
        <p>&emsp; For quick “at-a-glance” information the <a className='internal-link' href='/extension'>OSS-Score chrome extension </a> embeds scores directly into the GitHub repository homepage.</p>
        <p className='image-container'><img alt='extension-img'></img></p>
        <p>&emsp; For a deeper dive the <a className='internal-link' href='/'>OSS-Score website </a> allows users to compare the scores and individual metrics of different open source projects and tools.</p>
        <p className='image-container'><img alt='website-img'></img></p>
      </p>
      <header>Meet the Team</header>
      <div className='team-image-container'><img alt='OSS-Score Team Image'></img></div>
      <p className='meet-the-team'>
        <p className='team-blurb'>&emsp; OSS-Score was created by Vichka Fonarev, Eli Dow, Daniel Alderman, and Emil Polakiewicz as a senior computer science capstone project at Tufts University under the mentorship of Ole Lensmar.</p>  
          <div className='team-members'>
          <p className='team-member'>
            <div className='name'>Vichka Fonarev</div>
            <div className='photo'><img alt='vichka-pic'></img></div>
            <div className='blurb'>Okay I guess</div>
            <div className='linkedin'><a target="_blank" href='#'><FaLinkedin /></a></div>
          </p>

          <p className='team-member'>
            <div className='name'>Eli Dow</div>
            <div className='photo'><img alt='eli-pic'></img></div>
            <div className='blurb'>He's fine</div>
            <div className='linkedin'><a target="_blank" href='#'><FaLinkedin /></a></div>
          </p>

          <p className='team-member'>
            <div className='name'>Daniel Alderman</div>
            <div className='photo'><img alt='doniel-pic'></img></div>
            <div className='blurb'>Ew</div>
            <div className='linkedin'><a target="_blank" href='#'><FaLinkedin /></a></div>
          </p>

          <p className='team-member'>
            <div className='name'>Emil Polakiewicz</div>
            <div className='photo'><img alt='emil-pic'></img></div>
            <div className='blurb'>Pretty cool guy actually</div>
            <div className='linkedin'><a target="_blank" href='https://www.linkedin.com/in/emil-polakiewicz-12887b17b/'><FaLinkedin /></a></div>
          </p>
        </div>
      </p>
    </div>
  );
}

export default About;
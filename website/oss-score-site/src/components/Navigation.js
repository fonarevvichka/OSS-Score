import React from 'react'
import { Navbar, Container, Nav} from 'react-bootstrap'
import './Navigation.css';

import './Navigation.css';
import smallLogo from '../images/favicon.ico'
import gitLogo from '../images/Git-Icon-White.png'
import { ImGit } from "react-icons/im";

const Navigation = () => {

    const changeTabColor = (hover) => (event) => {
        if (hover === "hover") {
            event.target.style.color = "white"
        } else {
            event.target.style.color = "#FFFFFF8C"
        }
    }

    return (
            <Navbar className='navigation-bar' expand="lg" variant="dark">
                <Navbar.Brand href="/" style={{paddingLeft:"10px"}}>
                    <img className='navLogo' src={smallLogo}></img>
                    <div className='navLogoText'>OSS-Score</div>
                </Navbar.Brand>
                    <Nav className="navbar-nav ms-auto">
                        <Nav.Link class="nav-home" href="/" style={{color:"#FFFFFF8C"}}
                            onMouseEnter={changeTabColor("hover")} onMouseLeave={changeTabColor("")}>Home</Nav.Link>
                        <Nav.Link class="nav-link" href="/about" style={{color:"#FFFFFF8C"}}
                            onMouseEnter={changeTabColor("hover")} onMouseLeave={changeTabColor("")}>About</Nav.Link>
                        <Nav.Link class="nav-link" href="/extension" style={{color:"#FFFFFF8C"}}
                            onMouseEnter={changeTabColor("hover")} onMouseLeave={changeTabColor("")}>Extension</Nav.Link>
                        <Nav.Link class="nav-link" href="/generate-scores" style={{color:"#FFFFFF8C"}}
                            onMouseEnter={changeTabColor("hover")} onMouseLeave={changeTabColor("")}>How We Score</Nav.Link>
                        <Nav.Link class="nav-link" href="/privacypolicy" style={{color:"#FFFFFF8C"}}
                            onMouseEnter={changeTabColor("hover")} onMouseLeave={changeTabColor("")}>Privacy Policy</Nav.Link>
                        
                        <Nav.Link href="https://github.com/fonarevvichka/OSS-Score"
                            target="_blank" class="OSSScoreLink"
                            style={{color: "#FFFFFF8C" }}
                            onMouseEnter={changeTabColor("hover")} onMouseLeave={changeTabColor("")}>
                            <ImGit size={21}/> fonarevvichka/OSS-Score
                        </Nav.Link>
                    </Nav>
            </Navbar>
    );
}

export default Navigation
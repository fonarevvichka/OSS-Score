import React from 'react'
import { FaFish } from "react-icons/fa";
import { GiBubbles } from "react-icons/gi";
import './NotFound.css';

export default function NotFound() {
    return (
         <div className='wrapper'>
    <div className="privacy-body">
        <div className='bubbles'> <GiBubbles size={40}/> </div>
        <br /><p className="privacy">Go fish... this page does not exist</p>
        <div className='fish'> <FaFish size={40}/> </div>
    </div>
    </div>
    )
}

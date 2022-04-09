import React from 'react'
import { FaFish } from "react-icons/fa";
import { GiBubbles } from "react-icons/gi";
import './NotFound.css';

export default function NotFound() {
    return (
         <div className='wrapper'>
    <div className="not-found-body">
        <br /><p className="not-found-text">Go fish... this page does not exist <div className='bubbles'> <GiBubbles size={40} /> </div></p>
        <div className='fish'> <FaFish size={40} /> </div>
        </div>
    </div>
    )
}

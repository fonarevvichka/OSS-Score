import React from 'react'

import './Extension.css';

export default function Extension() {
  return (
    <div>
      <header>Download the OSS-Score Extension here!</header> 
      <div id="extension_descript">The chrome extension provides users with a comprehensive activity 
          and license score by embedding those scores directly into
          the GitHub repo homepage
      </div>
      <div class="text-center"><button id="download">Download</button></div>
    </div>
  )
}

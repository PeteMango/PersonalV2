import React from 'react';
import './App.css';
import Header from './components/Header/header'
import Body from './components/Body/body'
import Footer from './components/Footer/footer'

function App() {
  React.useEffect(() => {
      const event = new Event('rendered')
      document.dispatchEvent(event)
  }, [])

  return (
    <code>
      <div className="section_background-wrap">
        <canvas id="gradient-canvas" data-js-darken-top data-transition-in></canvas>
      </div>
      <div className='container'>
        <Header />
        <Body />
        <Footer />
      </div>
    </code>
  )
}

export default App;

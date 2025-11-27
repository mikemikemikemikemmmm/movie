import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import "./index.css"
import { LoadingComponent } from './components/loading'
import { SeatsContainer } from './components/seatContainer'
createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <div style={{ minHeight: "100vh", width: "100vw", position: "relative" }}>
      <LoadingComponent />
      <div style={{ width: 700, textAlign: "center", marginLeft: "auto", marginRight: "auto" }}>
        <SeatsContainer />
      </div>
    </div>
  </StrictMode>,
)

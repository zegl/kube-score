import styled from 'styled-components'
import kubeScoreLogo from './assets/logo.svg'
import Run from './Run'

function App() {
  return (
    <RootContainer>
      <div style={{ display: 'flex', gap: 16, alignItems: "center", justifyItems: "center", justifyContent: "space-between", flex: 1, flexShrink: 0, flexWrap: "wrap", width: "100%" }}>
        <div style={{ display: 'flex', gap: 8, alignItems: "center", flex: 1, flexShrink: 0, }}>
          <img src={kubeScoreLogo} style={{ height: 60 }} alt="kube-score logo" />
          <h1 style={{ fontSize: 32, textWrap: "nowrap" }}>kube-score</h1>
        </div>

        <div style={{ display: 'flex', gap: 16, alignItems: "center", fontSize: 16, flex: 1, flexShrink: 0 }}>
          <div style={{ flex: 1 }}></div>
          <a href="https://github.com/zegl/kube-score">GitHub</a>
          <a href="https://github.com/zegl/kube-score">README</a>
          <a href="https://github.com/zegl/kube-score/blob/master/README_CHECKS.md">Docs</a>
        </div>

      </div>

      <div style={{ fontSize: 14 }}>
        <p>Kubernetes object analysis with recommendations for improved reliability and security.</p>
        <p><code>kube-score</code> is a tool that does static code analysis of your Kubernetes object definitions. The
          output is
          a list of
          recommendations of what you can improve to make your application more secure and resilient.</p>
        <p><code>kube-score</code> is <a href="https://github.com/zegl/kube-score">open-source and available under the
          MIT-license</a>. For more information about how to use kube-score, see <a
            href="https://github.com/zegl/kube-score">zegl/kube-score</a> on GitHub.
          Use this website to easily test kube-score, just paste your object definition YAML or JSON in the box below.
        </p>

        <p>
          This tool is running 100% in your browser, no data is sent to any server!
        </p>
      </div>

      <Run />
    </RootContainer>
  )
}


const RootContainer = styled.div`
  margin: 0 auto;
  padding: 0.5;
  text-align: left;

  @media (min-width: 1024px) {
    padding: 2rem;
  }
`

export default App

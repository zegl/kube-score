import kubeScoreLogo from './assets/logo.svg'
import './App.css'
import Run from './Run'

function App() {
  return (
    <>
      <div style={{ display: 'flex', gap: 8, alignItems: "center" }}>
        <img src={kubeScoreLogo} style={{ height: 60 }} alt="kube-score logo" />
        <h1 style={{ fontSize: 32 }}>kube-score</h1>
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
    </>
  )
}

export default App

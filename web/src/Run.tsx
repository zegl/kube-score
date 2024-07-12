import { useEffect, useState } from "react";
import './assets/term.css'
import styled from "styled-components";

type Format = "human" | "ci" | "json" | "sarif" | "junit"

const FormatNames: Record<Format, string> = {
    "human": "CLI",
    "ci": "CI",
    "json": "JSON",
    "sarif": "SARIF",
    "junit": "JUnit",
}

// Types for the WebAssembly module
declare global {
    interface Window {
        handleScore: (input: string, type: Format) => string;
    }
}


function Run() {
    const [input, setInput] = useState(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: statefulset-test-1
spec:
  template:
    metadata:
      labels:
        app: foo
    spec:
      containers:
      - name: foobar
        image: foo:bar
---
apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: app-budget
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app: not-foo`)


    const [wasmLoading, setWasmLoading] = useState(false)
    const [result, setResult] = useState("")
    const [format, setFormat] = useState<Format>("human")

    const scoreWhenReady = () => {
        if (typeof window.handleScore === "undefined") {
            setWasmLoading(true)
            setTimeout(scoreWhenReady, 100);
            return;
        }

        score();
    }

    const score = () => {
        setWasmLoading(false)
        setResult("Running...")
        setResult(window.handleScore(input, format));
    }

    useEffect(() => {
        scoreWhenReady()
    }, [input, format])

    const output = wasmLoading ? "LOADING_WASM" :
        !input ? "NO_INPUT" :
            result ? "RESULT" : "RUNNING" as const

    return (
        <div>
            <div style={{color: "#4e4e4e"}}>Output format</div>
                <div style={{ display: "flex", gap: 2, alignItems: "center" }}>
                    {Object.entries(FormatNames).map(([f, name]) => <OutputFormatButton $selected={format == f} key={f} onClick={() => setFormat(f as Format)}>{name}</OutputFormatButton>)}
                </div>


        <InputOutputContainer>
            <InputTextarea value={input} onChange={(e) => setInput(e.target.value)}></InputTextarea>

            <div style={{ flex: 1, flexShrink: 0, overflow: "hidden", display: "flex", gap: 4, flexDirection: "column" }}>
            
                

                <TermContainer>
                    {output === "LOADING_WASM" && <div>Loading WebAssembly...</div>}
                    {output === "NO_INPUT" && <div>Paste your Kubernetes YAML in the textarea to get started</div>}
                    {output === "RESULT" && <>
                        {format === "human" ? <div dangerouslySetInnerHTML={{ __html: result }}></div> : <div>{result}</div>}
                    </>}
                    {output === "RUNNING" && <div>Running...</div>}
                </TermContainer></div>
        </InputOutputContainer>
        </div>
    );
}

const InputOutputContainer = styled.div`
    display: flex;
    gap: 8px;
    margin-top: 16px;

    flex-direction: column;

    @media (min-width: 1024px) {
        flex-direction: row;
  }
`

const OutputFormatButton = styled.button<{ $selected: boolean; }>`
    background: ${props => props.$selected ? "#c172a9" : "#474747"};
    color: ${props => props.$selected ? "white" : "white"};
    border-radius: 4px;
`

const TermContainer = styled.div`
    background: #171717;
    border-radius: 5px;
    color: white;
    word-break: break-word;
    overflow-wrap: break-word;
    font-family: "SFMono-Regular", Monaco, Menlo, Consolas, "Liberation Mono", Courier, monospace;
    font-size: 14px;
    line-height: 20px;
    padding: 14px 18px;
    white-space: pre-wrap;
    flex: 1;
    flex-shrink: 0;
`

const InputTextarea = styled.textarea`
    border-radius: 5px;
    font-family: "SFMono-Regular", Monaco, Menlo, Consolas, "Liberation Mono", Courier, monospace;
    font-size: 14px;
    line-height: 20px;
    padding: 14px 18px;
    flex: 1;
    flex-shrink: 0;
    min-height: 200px
`

export default Run;
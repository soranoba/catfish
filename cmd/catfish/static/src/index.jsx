import * as React from "react";
import {createRoot} from 'react-dom/client';
import {Accordion, Button, Table} from "react-bootstrap";

let apiEndpoint = __API_ENDPOINT__ || "";

async function fetchConfig() {
  return fetch(apiEndpoint + "/api/config").then((res) => {
    if (!res.ok) {
      throw new Error('Failed to load config');
    }
    return res.json();
  });
}

async function fetchVariables() {
  return fetch(apiEndpoint + "/api/variables").then((res) => {
    if (!res.ok) {
      throw new Error('Failed to load variables');
    }
    return res.json();
  });
}

const App = () => {
  const [config, setConfig] = React.useState({routes: []});
  const [variables, setVariables] = React.useState({global_variables: {}, route_variables: []});

  React.useEffect(() => {
    fetchConfig().then(setConfig);
    fetchVariables().then(setVariables);
  }, []);

  async function resetVariables() {
    await fetch(apiEndpoint + "/api/variables/reset", {
      method: "PUT",
    }).then((res) => {
      if (!res.ok) {
        throw new Error('Failed to reset variables');
      }
    });
    fetchVariables().then(setVariables);
  }

  return (
    <>
      <nav className="navbar navbar-light bg-light mb-4">
        <div className="container-fluid">
          <span className="navbar-brand mb-0 h1">Catfish</span>
        </div>
      </nav>
      <div className="container">
        <h2>Config</h2>
        <Accordion className="mb-4">
          {(config.routes || []).map((route) => {
            const key = `${route.method}-${route.path}`;
            return (
              <Accordion.Item eventKey={key} key={key}>
                <Accordion.Header>
                  <span style={{width: 75, fontWeight: "bold"}}>{route.method}</span>
                  <span>{route.path}</span>
                </Accordion.Header>
                <Accordion.Body>
                  {route.response.map((res) => (
                    <div key={res.name}>
                      <h5>{res.name}</h5>
                      <Table bordered hover className="mb-4">
                        <tbody>
                        <tr>
                          <th>status</th>
                          <td>{res.status}</td>
                        </tr>
                        <tr>
                          <th>cond</th>
                          <td>{res.cond}</td>
                        </tr>
                        <tr>
                          <th>delay</th>
                          <td>{res.delay}</td>
                        </tr>
                        <tr>
                          <th>header</th>
                          <td>{JSON.stringify(res.header)}</td>
                        </tr>
                        <tr>
                          <th>body</th>
                          <td>{res.body}</td>
                        </tr>
                        </tbody>
                      </Table>
                    </div>
                  ))}
                </Accordion.Body>
              </Accordion.Item>
            );
          })}
        </Accordion>
        <h2>Variables</h2>
        <div className="mb-4">
          <Button variant="danger" size="sm" onClick={resetVariables}>Rest all variables</Button>
        </div>
        <h4>Global</h4>
        <Table bordered hover className="mb-4">
          <tbody>
          {Object.keys(variables.global_variables).map((key) => (
            <tr key={key}>
              <th>{key}</th>
              <td>{variables.global_variables[key]}</td>
            </tr>
          ))}
          </tbody>
        </Table>
        <h4>Route</h4>
        <Accordion className="mb-4">
          {(variables.route_variables || []).map((route_variable) => {
            const key = `${route_variable.route.method}-${route_variable.route.path}`;
            return (
              <Accordion.Item eventKey={key} key={key}>
                <Accordion.Header>
                  <span style={{width: 75, fontWeight: "bold"}}>{route_variable.route.method}</span>
                  <span>{route_variable.route.path}</span>
                </Accordion.Header>
                <Accordion.Body>
                  <Table bordered hover>
                    <tbody>
                    {Object.keys(route_variable.variables).map((key) => (
                      <tr key={key}>
                        <th>{key}</th>
                        <td>{route_variable.variables[key]}</td>
                      </tr>
                    ))}
                    </tbody>
                  </Table>
                </Accordion.Body>
              </Accordion.Item>
            );
          })}
        </Accordion>
      </div>
    </>
  );
};

const root = createRoot(document.getElementById("root"));
root.render(<App/>);

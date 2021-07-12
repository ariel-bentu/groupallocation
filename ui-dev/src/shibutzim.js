import React, { useEffect, useState } from 'react';
import { List, ListItem, ListItemText, Paper, Button } from '@material-ui/core';
import useStyles from "./styles.js"
import * as api from './api'
import { EditResults } from './edit-result';
import { VBox, HBox, Spacer, Header } from './elems'

function Shibutzim(props) {
  const classes = useStyles();
  const [results, setResults] = useState([]);
  const [currentResult, setCurrentResult] = useState("");
  const [editResults, setEditResults] = useState(false);

  const loadResults = async () => {
    return api.loadResults(props.currentTask).then(results => {
      console.log("load results")
      setResults(results);
    })
  };

  useEffect(() => {
    setCurrentResult("")
    loadResults()
  }, [props.currentTask]);// eslint-disable-line react-hooks/exhaustive-deps

  //console.log(currentResult ? results.find(r => r.id === currentResult).name : "")
  return (
    <div className={classes.paperContainer}>
      <Paper elevation={3} className={classes.paper}>
        <Header>שיבוצים</Header>
        <VBox>
          <List className={classes.list} style={{ height: 250 }}>
            {props.items ? props.items.map((item) => (
              <ListItem className={classes.listItem}
                button selected={props.currentTask === item.id} onClick={() => props.onChangeTask(item.id)}>
                <ListItemText primary={item.name} />
              </ListItem>
            )) : null}
          </List>
          <Spacer />
          <input type="button" id="btnDeleteTask" value="מחק שיבוץ" />
        </VBox>
      </Paper>



      <Paper elevation={3} className={classes.paper}>
        <Header>הרצה</Header>
        <VBox>
          <Spacer />
          limit
          <input type="input" id="runLimit" value="20" /><br />
          grace-level
          <input type="input" id="graceLevel" value="0" /><br />
          sensitive-to-only-last
          <input type="input" id="sensitiveToOnlyLast" value="0" />
          <Spacer />
          <Button
            variant="outlined"
            color="primary"
            onClick={() => {
              let limit = 20 //$("#runLimit").val()
              let graceLevel = 0 //$("#graceLevel").val()
              let sensitiveToOnlyLast = 0 //$("#sensitiveToOnlyLast").val()
              window.open("/run?task=" + props.currentTask + "&limit=" + limit + "&graceLevel=" + graceLevel + "&sensitiveToOnlyLast=" + sensitiveToOnlyLast);
            }
            }
          >הרץ...</Button>
        </VBox>
      </Paper>



      <Paper elevation={3} >
        <Header>תוצאות</Header>
        <VBox>
          <List className={classes.list} style={{ height: 250 }}>
            {results.map((item) => (
              <ListItem button className={classes.listItem} selected={currentResult === item.id} onClick={() => setCurrentResult(item.id)}>
                <ListItemText primary={item.title !== "" ? item.title : item.runDate} />
              </ListItem>
            ))}
          </List>
          <Spacer />
          <HBox>
            <input type="button" id="btnShowResult" value="הצג תוצאה" />
            <input type="button" id="btnShowCleanResult" value="הצג תוצאה - נקי" />
            <input type="button" id="btnDeleteResult" value="מחק תוצאה" />
            <input type="button" id="btnDuplicateResult" value="שכפל תוצאה" />

          </HBox>
          <HBox>
            <Button
              variant="outlined"
              color="primary"
              onClick={() => setEditResults(true)}
              disabled={currentResult === ""}
            >ערוך...</Button>
            <EditResults
              open={editResults}
              Name={currentResult !== "" ? results.find(r => r.id === currentResult).title : ""}
              Cancel={() => setEditResults(false)}
              Save={(newName) => {
                api.saveResultName(props.currentTask, currentResult, newName).then(() => loadResults())
                setEditResults(false);
              }}
            />

          </HBox>
        </VBox>
      </Paper>
    </div>
  );
}

export default Shibutzim;
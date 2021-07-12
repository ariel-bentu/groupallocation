import React, { useEffect, useState } from 'react';
import { List, ListItem, ListItemText, Paper, TextField } from '@material-ui/core';
import useStyles from "./styles.js"
import * as api from './api'
import { EditResults } from './edit-result';
import { VBox, HBox, Spacer, Header, GButton } from './elems'

function Shibutzim(props) {
  const classes = useStyles();
  const [results, setResults] = useState([]);
  const [currentResult, setCurrentResult] = useState("");
  const [editResults, setEditResults] = useState(false);
  const [timeLimit, setTimeLimit] = useState(20);
  const [graceLevel, setGraceLevel] = useState(0);

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
                <ListItemText className={classes.listItemText} primary={item.name} />
              </ListItem>
            )) : null}
          </List>
          <Spacer />
          <GButton label="מחק שיבוץ" onClick={() => alert("todo")} />
        </VBox>
      </Paper>



      <Paper elevation={3} className={classes.paper}>
        <Header>הרצה</Header>
        <VBox>
          <Spacer />
          <TextField label="מגבלת זמן" value={timeLimit} onChange={(e) => setTimeLimit(e.currentTarget.value)} />
          <Spacer />
          <TextField label="רגישות" value={graceLevel} onChange={(e) => setGraceLevel(e.currentTarget.value)} />
          <Spacer />
          <GButton label="הרץ..." onClick={() => {
            let limit = timeLimit
            let graceLevel = graceLevel
            let sensitiveToOnlyLast = 0 ;
            window.open("/run?task=" + props.currentTask + "&limit=" + limit + "&graceLevel=" + graceLevel + "&sensitiveToOnlyLast=" + sensitiveToOnlyLast);
          }} />
        </VBox>
      </Paper>



      <Paper elevation={3} >
        <Header>תוצאות</Header>
        <VBox>
          <List className={classes.list} style={{ height: 250 }}>
            {results.map((item) => (
              <ListItem button className={classes.listItem} selected={currentResult === item.id} onClick={() => setCurrentResult(item.id)}>
                <ListItemText className={classes.listItemText} primary={item.title !== "" ? item.title : item.runDate} />
              </ListItem>
            ))}
          </List>
          <Spacer />
          <HBox>
            <GButton label="שנה שם..." onClick={() => setEditResults(true)} disabled={currentResult === ""} />
            <GButton label="הצג תוצאה" disabled={!currentResult} onClick={() => alert("todo")} />
            <GButton label="הצג תוצאה - נקי" disabled={!currentResult} onClick={() => alert("todo")} />
            <GButton label="מחק תוצאה" disabled={!currentResult} onClick={() => alert("todo")} />
            <GButton label="שכפל תוצאה" disabled={!currentResult} onClick={() => alert("todo")} />
          </HBox>
          <HBox>
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
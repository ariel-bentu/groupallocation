import React, { useEffect, useState } from 'react';
import { List, ListItem, ListItemText, TextField } from '@material-ui/core';
import useStyles from "./styles.js"
import * as api from './api'
import { EditResults } from './edit-result';
import EditTask from './edit-task.js';
import { VBox, HBox, Spacer, Header, GButton, Paper1 } from './elems'

function Shibutzim(props) {
  const classes = useStyles();
  const [results, setResults] = useState([]);
  const [currentResult, setCurrentResult] = useState("");
  const [editResults, setEditResults] = useState(false);
  const [timeLimit, setTimeLimit] = useState(20);
  const [graceLevel, setGraceLevel] = useState(0);
  const [editTaskDialog, setEditTaskDialog] = useState(undefined);


  const loadResults = async () => {
    return api.loadResults(props.currentTask).then(results => {
      console.log("load results")
      setResults(results);
    })
  };

  const getTaskName = (id) => {
    let res = props.items.find(r => r.id === id);
    return res && res.name !== "" ? res.name : "";
  }
  const getResultName = (id) => {
    let res = results.find(r => r.id === id);
    if (!res) return ""
    return res.title !== "" ? res.title : res.runDate;
  }

  useEffect(() => {
    setCurrentResult("")
    loadResults()
  }, [props.currentTask]);// eslint-disable-line react-hooks/exhaustive-deps

  return (
    <div className={classes.paperContainer}>
      <Paper1 elevation={3} className={classes.paper}>
        <Header>שיבוצים</Header>
        <VBox>
          <List className={classes.list} style={{ height: 250 }}>
            {props.items ? props.items.map((item) => (
              <ListItem className={classes.listItem} key={item.id}
                button selected={props.currentTask === item.id} onClick={() => props.onChangeTask(item.id)}>
                <ListItemText className={classes.listItemText} primary={item.name} />
              </ListItem>
            )) : null}
          </List>
          <Spacer />
          <GButton label="שיבוץ חדש..." onClick={() => setEditTaskDialog({ id: "" })} />
          <GButton label="מחק שיבוץ" disabled={!props.currentTask} onClick={() => props.msg.alert({
            title: "מחיקת שיבוץ",
            message: `האם למחוק את השיבוץ '${getTaskName(props.currentTask)}' \nמחיקת השיבוץ הינה בלתי הפיכה!!!`,
            buttons: [{
              label: "מחק",
              callback: () => {
                api.deleteTask(props.currentTask).then(() => {
                  props.reloadTasks(undefined);
                })
              }
            },
            {
              label: "בטל",
              callback: () => { }
            }]
          })} />
          <EditTask open={editTaskDialog}
            task={editTaskDialog}
            Name={getTaskName(props.currentTask)}
            Cancel={() => setEditTaskDialog(undefined)}
            Save={(newTask) => {
              api.saveTask(newTask).then(() => props.reloadTasks(newTask.id));
              setEditTaskDialog(undefined);
            }}
          />
        </VBox>

      </Paper1>




      <Paper1 elevation={3} className={classes.paper}>
        <Header>הרצה</Header>
        <VBox>
          <Spacer />
          <TextField label="מגבלת זמן" value={timeLimit} onChange={(e) => setTimeLimit(e.currentTarget.value)} />
          <Spacer />
          <TextField label="רגישות" value={graceLevel} onChange={(e) => setGraceLevel(e.currentTarget.value)} />
          <Spacer />
          <GButton label="הרץ..." onClick={() => {
            let sensitiveToOnlyLast = 0;
            window.open("/run?task=" + props.currentTask + "&limit=" + timeLimit + "&graceLevel=" + graceLevel + "&sensitiveToOnlyLast=" + sensitiveToOnlyLast);
          }} />

        </VBox>
      </Paper1>



      <Paper1 elevation={3} >
        <Header>תוצאות</Header>
        <VBox>
          <List className={classes.list} style={{ height: 250 }}>
            {results.map((item) => (
              <ListItem button className={classes.listItem} selected={currentResult === item.id}
                onClick={() => setCurrentResult(item.id)} key={item.id}>
                <ListItemText className={classes.listItemText} primary={item.title !== "" ? item.title : item.runDate} />
              </ListItem>
            ))}
          </List>
          <Spacer />
          <HBox>
            <GButton label="שנה שם..." onClick={() => setEditResults(true)} disabled={currentResult === ""} />
            <GButton label="הצג תוצאה" disabled={!currentResult} onClick={() => api.showResults(props.currentTask, currentResult, false)} />
            <GButton label="הצג תוצאה - נקי" disabled={!currentResult} onClick={() => api.showResults(props.currentTask, currentResult, false)} />
            <GButton label="מחק תוצאה" disabled={!currentResult} onClick={() => props.msg.alert({
              title: "מחיקת תוצאה",
              message: `האם למחוק את התוצאה '${getResultName(currentResult)}'`,
              buttons: [{
                label: "מחק",
                callback: () => {
                  api.deleteResult(currentResult.id).then(() => {
                    setCurrentResult(undefined);
                    api.loadResults(props.currentTask).then(rst => setResults(rst));
                  })
                }
              },
              {
                label: "בטל",
                callback: () => { }
              }]
            })} />
            <GButton label="שכפל תוצאה" disabled={!currentResult} onClick={() =>
              api.duplicateResult(props.currentTask, currentResult.id, "העתק של " + currentResult.name).then(() => {
                api.loadResults(props.currentTask).then(rst => setResults(rst));
              })} />
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
      </Paper1>
    </div>
  );
}

export default Shibutzim;
import React from "react";
import ReactDOM from "react-dom";

import "./styles.css";
import useStyles from "./styles.js"
import PropTypes from "prop-types";
import { AppBar, Tabs, Tab, Typography, Box, Snackbar } from "@material-ui/core"
import { Alert } from './elems'

import Shibutzim from './shibutzim.js'
import PupilPrefs from './pupil-prefs'
import SubGroups from './sub-groups'
import * as api from './api.js'
import AlertDialog from './alert-dialog'

function TabPanel(props) {
  const { children, value, index, ...other } = props;

  return (
    <Typography
      component="div"
      role="tabpanel"
      hidden={value !== index}
      id={`scrollable-auto-tabpanel-${index}`}
      {...other}
    >
      {props.title ? <h1 align="center">{props.title}</h1> : null}
      <Box className={props.class} p={3}>{children}</Box>
    </Typography>
  );
}

TabPanel.propTypes = {
  children: PropTypes.node,
  index: PropTypes.any.isRequired,
  value: PropTypes.any.isRequired
};




export default function App() {
  const classes = useStyles();
  const [value, setValue] = React.useState(0);
  const [listTasks, setListTasks] = React.useState([]);
  const [listPupils, setListPupils] = React.useState([]);
  const [currentTask, setCurrentTask] = React.useState("");
  const [message, setMessage] = React.useState(undefined);
  const [alert, setAlert] = React.useState(undefined);
  const [tabDirty, setTabDirty] = React.useState(false);

  React.useEffect(() => {
    console.log("load tasks")
    api.loadTasks().then(tasks => {
      setListTasks(tasks);
      if (tasks.length > 0)
        setCurrentTask(tasks[0].id)
      else
        setCurrentTask("")
    })
  }, []);

  const msg = {
    alert: (al) => setAlert(al),
    notify: (msg) => setMessage(msg)
  }

  const _setTabDirty = (val) => setTabDirty(val);

  React.useEffect(() => {
    console.log("load pupils")

    api.loadPupils(currentTask).then(pupils => {
      setListPupils(pupils);
    })
  }, [currentTask]);


  function changeTab(event, newValue) {
    if (tabDirty) {
      msg.alert({
        title: "שינויים לא נשמרו",
        message: "לפני החלפת לשונית יש לשמור או לבטל השינויים",
        buttons: [
          {
            label: "בטל",
            callback: () => { }
          }
        ]
      })
    } else
      setValue(newValue);
  }

  return (
    <div>
      <AppBar position="static" color="default" >
        <Tabs
          value={value}
          onChange={changeTab}
          indicatorColor="primary"
          textColor="primary"
          variant="fullwidth"
          scrollButtons="auto"
          centered
        >
          <Tab label={"שיבוצים" + (value === 0 && tabDirty ? "*" : "")}  />
          <Tab label={"תלמידים וחברים" + (value === 2 && tabDirty ? "*" : "")}  />
          <Tab label={"קבוצות" + (value === 3 && tabDirty ? "*" : "")} />
        </Tabs>
      </AppBar>
      <TabPanel value={value} index={0} class={classes.root}>
        <Shibutzim items={listTasks} currentTask={currentTask} onChangeTask={(id) => setCurrentTask(id)} msg={msg} setDirty={_setTabDirty} />
      </TabPanel>
      <TabPanel value={value} index={1} class={classes.root}>
           <PupilPrefs currentTask={currentTask} pupils={listPupils} msg={msg} setDirty={_setTabDirty} />
      </TabPanel>
      <TabPanel value={value} index={2} class={classes.root}>
        <SubGroups currentTask={currentTask} pupils={listPupils} msg={msg} setDirty={_setTabDirty} />
      </TabPanel>


      <Snackbar open={message} autoHideDuration={6000} onClose={(event, reason) => (reason === 'clickaway') ? {} : setMessage(undefined)}>
        <Alert onClose={(event, reason) => (reason === 'clickaway') ? {} : setMessage(undefined)}
          severity="success">
          {message}
        </Alert>
      </Snackbar>
      <AlertDialog open={alert} alert={alert} close={() => setAlert(undefined)} />
    </div>
  );
}

const rootElement = document.getElementById("root");
ReactDOM.render(<App />, rootElement);

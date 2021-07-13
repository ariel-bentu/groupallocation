import React, { useEffect, useState } from 'react';

import {
    Dialog, DialogTitle, DialogContent,
    DialogActions, 
    Slide, TextField
} from '@material-ui/core';

import { VBox, LabeledCheckBox, GButton } from './elems';

const Transition = React.forwardRef(function Transition(props, ref) {
    return <Slide direction="up" ref={ref} {...props} />;
});


export default function EditTask(props) {
    const [editTask, setEditTask] = useState(undefined);

    useEffect(() => {
        console.log("set editTask")
        setEditTask(props.task);
    }, [props.open]) // eslint-disable-line react-hooks/exhaustive-deps

    const toggle = (propName) => {
        let newTask = { ...editTask }
        newTask[propName] = !newTask[propName];
        setEditTask(newTask);
    }
    const onTextChange = (e, propName) => {
        let newTask = { ...editTask }
        newTask[propName] = e.currentTarget.value
        setEditTask(newTask);
    }

    return (
        editTask ? <Dialog
            open={props.open}
            TransitionComponent={Transition}
            keepMounted
            onClose={props.Cancel}
        >
            <DialogTitle>{editTask && editTask.id !== undefined ? "ערוך פרטי שיבוץ" : "שיבוץ חדש"}</DialogTitle>
            <DialogContent>
                <VBox style={{alignItems:'flex-end'}}>
                    <TextField
                        onChange={(e) => {
                            onTextChange(e, "name")
                        }}
                        inputProps={{ style: { textAlign: 'right' } }}
                        label="שם"
                        variant="filled"
                        value={editTask.name}
                    />
                    <TextField
                        onChange={(e) => onTextChange(e, "numOfGroups")}
                        inputProps={{ style: { textAlign: 'right' } }}
                        label="מספר כיתות"
                        variant="filled"
                        value={editTask.numOfGroups}
                    />
                    <LabeledCheckBox
                        checked={editTask.encrypted}
                        onClick={(e) => toggle("encrypted")}
                        label={"מוצפן"}
                    />
                </VBox>

            </DialogContent>
            <DialogActions>
                <GButton onClick={props.Cancel} label="בטל"/>
                <GButton onClick={() => props.Save(editTask)} label="שמור"/>
            </DialogActions>
        </Dialog> : null)
}
import React, { useEffect, useState } from 'react';

import {
    Dialog, DialogTitle, DialogContent,
    DialogActions, Button,
    Slide, TextField
} from '@material-ui/core';

import { VBox, LabeledCheckBox } from './elems';

const Transition = React.forwardRef(function Transition(props, ref) {
    return <Slide direction="up" ref={ref} {...props} />;
});


export default function EditGroup(props) {
    const [editGroup, setEditGroup] = useState(undefined);

    useEffect(() => {
        console.log("set EditGroup")
        setEditGroup(props.group);
    }, [props.open]) // eslint-disable-line react-hooks/exhaustive-deps

    const toggle = (propName) => {
        let newGroup = { ...editGroup }
        newGroup[propName] = !newGroup[propName];
        setEditGroup(newGroup);
    }

    return (
        editGroup ? <Dialog
            open={props.open}
            TransitionComponent={Transition}
            keepMounted
            onClose={props.Cancel}
        >
            <DialogTitle>{editGroup && editGroup.id !== undefined ? "ערוך קבוצה" : "קבוצה חדשה"}</DialogTitle>
            <DialogContent>
                <VBox style={{alignItems:'flex-end'}}>
                    <TextField
                        onChange={(e) => {
                            setEditGroup({ ...editGroup, name: e.currentTarget.value })
                        }}
                        inputProps={{ style: { textAlign: 'right' } }}
                        label="שם"
                        variant="filled"
                        value={editGroup.name}
                    />
                    <LabeledCheckBox
                        checked={editGroup.isGarden}
                        onClick={(e) => toggle("isGarden")}
                        label={"קבוצת גן"}
                    />
                    <LabeledCheckBox
                        checked={editGroup.isUnite}
                        onClick={(e) => toggle("isUnite")}
                        label={"קבוצת איחוד"}
                    />
                    <LabeledCheckBox
                        checked={editGroup.isSpreadEvenly}
                        onClick={(e) => toggle("isSpreadEvenly")}
                        label={"פיזור אחיד"}
                    />
                    <LabeledCheckBox
                        checked={editGroup.isGenderSensitive}
                        onClick={(e) => toggle("isGenderSensitive")}
                        label={"רגיש לבנים/בנות"}
                    />
                    <LabeledCheckBox
                        checked={editGroup.isInactive}
                        onClick={(e) => toggle("isInactive")}
                        label={"פעיל"}
                    />
                   
                </VBox>

            </DialogContent>
            <DialogActions>
                <Button onClick={props.Cancel} color="primary">
                    בטל
                </Button>
                <Button onClick={() => props.Save(editGroup)} color="primary">
                    שמור
                </Button>
            </DialogActions>
        </Dialog> : null)
}
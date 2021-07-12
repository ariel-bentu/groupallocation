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


export default function EditPupil(props) {
    const [editPupil, setEditPupil] = useState(undefined);

    useEffect(() => {
        console.log("set EditPupil")
        setEditPupil(props.pupil);
    }, [props.open]) // eslint-disable-line react-hooks/exhaustive-deps

    const toggle = (propName) => {
        let newPupil = { ...editPupil }
        newPupil[propName] = !newPupil[propName];
        setEditPupil(newPupil);
    }

    const onTextChange = (e, propName) => {
        let newPupil = { ...editPupil }
        newPupil[propName] = e.currentTarget.value
        setEditPupil(newPupil);
    }

    return (
        editPupil ? <Dialog
            open={props.open}
            TransitionComponent={Transition}
            keepMounted
            onClose={props.Cancel}
        >
            <DialogTitle>{editPupil && editPupil.id !== undefined ? "ערוך פרטי תלמיד" : "הוספת תלמיד/ה"}</DialogTitle>
            <DialogContent>
                <VBox style={{alignItems:'flex-end'}}>
                    <TextField
                        onChange={(e) => {
                            onTextChange(e, "name")  
                        }}
                        inputProps={{ style: { textAlign: 'right' } }}
                        label="שם"
                        variant="filled"
                        value={editPupil.name}
                    />
                    <LabeledCheckBox
                        checked={editPupil.isMale}
                        onClick={(e) => toggle("isMale")}
                        label={"בן"}
                    />
                    <LabeledCheckBox
                        checked={editPupil.active}
                        onClick={(e) => toggle("active")}
                        label={"פעיל"}
                    />
                    <TextField
                        onChange={(e) => {
                            onTextChange(e, "remarks")  
                        }}
                        inputProps={{ style: { textAlign: 'right' } }}
                        label="הערות"
                        variant="filled"
                        value={editPupil.remarks}
                    />
                   
                </VBox>

            </DialogContent>
            <DialogActions>
                <Button onClick={props.Cancel} color="primary">
                    בטל
                </Button>
                <Button onClick={() => props.Save(editPupil)} color="primary">
                    שמור
                </Button>
            </DialogActions>
        </Dialog> : null)
}
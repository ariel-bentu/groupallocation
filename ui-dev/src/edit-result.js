import React ,{useEffect, useState} from 'react';

import { Dialog, DialogTitle, DialogContent, 
     DialogActions, Button,
    Slide, TextField } from '@material-ui/core';

const Transition = React.forwardRef(function Transition(props, ref) {
  return <Slide direction="up" ref={ref} {...props} />;
});


export function EditResults(props) {
    const [name, setName] = useState("");

    useEffect(()=>{
        console.log("setName")
        setName(props.Name);
    }, [props.open])

    return (
    <Dialog
        open={props.open}
        TransitionComponent={Transition}
        keepMounted
        onClose={props.Cancel}
        aria-labelledby="alert-dialog-slide-title"
        aria-describedby="alert-dialog-slide-description"
      >
        <DialogTitle id="alert-dialog-slide-title">{"שנה שם של תוצאה"}</DialogTitle>
        <DialogContent>
          <TextField 
            onChange={(e)=>setName(e.target.value)} 
            inputProps={{style: { textAlign: 'right' }}} 
            label="שם התוצאה" 
            variant="filled" 
            value={name}
            helperText="הכנס שם חדש"
        />
        </DialogContent>
        <DialogActions>
          <Button onClick={props.Cancel} color="primary">
            בטל
          </Button>
          <Button onClick={()=>props.Save(name)} color="primary">
            שמור
          </Button>
        </DialogActions>
      </Dialog>)
}
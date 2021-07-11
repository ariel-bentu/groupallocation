import React  from 'react';

import { Dialog, DialogTitle, DialogContent, 
     DialogActions, Button,
    Slide } from '@material-ui/core';

const Transition = React.forwardRef(function Transition(props, ref) {
  return <Slide direction="up" ref={ref} {...props} />;
});


export default function AlertDialog(props) {

    return (
        props.alert?
    <Dialog
        open={props.open}
        TransitionComponent={Transition}
        keepMounted
      >
        <DialogTitle>{props.alert.title}</DialogTitle>
        <DialogContent>{props.alert.message}</DialogContent>
        <DialogActions>
        {props.alert.buttons.map(btn=>(
          <Button onClick={()=>{
                  btn.callback()
                  props.close()
              }} color="primary">{btn.label}</Button>))
        }
        
        </DialogActions>
      </Dialog>:null)
}
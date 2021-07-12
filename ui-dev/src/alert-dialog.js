import React from 'react';

import {
    Dialog, DialogTitle, DialogContent,
    DialogActions, Slide
} from '@material-ui/core';
import { GButton } from './elems'

const Transition = React.forwardRef(function Transition(props, ref) {
    return <Slide direction="up" ref={ref} {...props} />;
});


export default function AlertDialog(props) {

    return (
        props.alert ?
            <Dialog
                open={props.open}
                TransitionComponent={Transition}
                keepMounted
            >
                <DialogTitle>{props.alert.title}</DialogTitle>
                <DialogContent>{props.alert.message}</DialogContent>
                <DialogActions>
                    {props.alert.buttons.map(btn => (
                        <GButton label={btn.label} onClick={() => {
                            btn.callback()
                            props.close()
                        }} />))
                    }

                </DialogActions>
            </Dialog> : null)
}
import { makeStyles } from "@material-ui/core/styles";

const useStyles = makeStyles(theme => ({
    root: {
        flexGrow: 1,
        width: "100%",
        backgroundColor: theme.palette.background.paper
    },
    tabs: {
        alignContent: "flex-start"
    },
    paperContainer: {
        direction:'rtl',
        display: 'flex',
            flexWrap: 'wrap',
        '& > *': {
            margin: theme.spacing(1),
            width: theme.spacing(55),
            height: theme.spacing(65),
        },
        backgroundColor: theme.palette.background.paper,
        justifyContent:'center'
    },
    paper : {
        
        justifyContent:'center'
    },
    list: {
        margin:5,
        width: '80%',
        overflow: 'auto',
        backgroundColor:'#E1E1E1'
    },
    listItem: {
        height: 30,
        alignItems:'center'
    },
    listItemText: {
        textAlign:'right'
    }
}));

export default useStyles;
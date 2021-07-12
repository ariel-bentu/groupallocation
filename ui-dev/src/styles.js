import { makeStyles } from "@material-ui/core/styles";

const useStyles = makeStyles(theme => ({
    root: {
        height:'100%',
        width: '98%',
        
        backgroundColor: 'gray'// theme.palette.background.paper
    },
    tabs: {
        alignContent: "flex-start"
    },
    paperContainer: {
        direction:'rtl',
        display: 'flex',
        flexDirection:'row',
        flexWrap: 'wrap',
        alignItems: 'stretch',
        '& > *': {
            margin: theme.spacing(1),
            width: theme.spacing(55),
            //height: '100%'//theme.spacing(65),
           
        },
        backgroundColor: theme.palette.background.paper,
        justifyContent:'center',
        width:'100%',
        height:'100%'
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
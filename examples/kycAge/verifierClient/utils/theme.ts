import {funk} from '@theme-ui/presets'

export default {
    ...funk, 
    colors: {
    ...funk.colors,
    polygon:"#8247e5"
    },
    layout:{
        allCenter: {
           justifyContent:"center",
           alignItems:"center" 
        }
    },
    text:{
       para: {
        textAlign: "center", fontSize: [24], mt: [3] 
       } 
    }
}
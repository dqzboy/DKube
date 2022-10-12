import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ELIcons from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'
import { GlobalCmComponent } from "codemirror-editor-vue3";
import 'codemirror/theme/idea.css'
import 'codemirror/mode/yaml/yaml.js'

const app = createApp(App)
for (let iconName in ELIcons) {
    app.component(iconName, ELIcons[iconName])
}
app.use(ElementPlus)
app.use(GlobalCmComponent, { componentName: "codemirror" });
app.use(router)
app.mount('#app')
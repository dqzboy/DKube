<template>
    <div class="common-layout">
        <el-container style="height: 100vh;">
            <el-aside class="aside" :width="asideWidth">
                <el-affix class="aside-affix" :z-index="1200">
                    <div class="aside-logo" >
                        <el-image class="logo-image" :src="logo" />
                        <span :class="[isCollapse ? 'is-collapse' : '']">
                            <span class="logo-name" >DKube</span>
                        </span>
                    </div>
                </el-affix>
                <el-menu class="aside-menu"
                    router
                    :default-active="$route.path"
                    :collapse="isCollapse"
                    background-color="#131b27"
                    text-color="#bfcbd9"
                    active-text-color="#20a0ff">
                    <div v-for="menu in routers" :key="menu">
                        <el-menu-item class="aside-menu-item" v-if="menu.children && menu.children.length == 1" :index="menu.children[0].path">
                        <el-icon><component :is="menu.children[0].icon" /></el-icon>
                        <template #title>
                            {{menu.children[0].name}}
                        </template>
                        </el-menu-item>
                        <el-sub-menu class="aside-submenu" v-else-if="menu.children" :index="menu.path">
                            <template #title>
                                <el-icon><component :is="menu.icon" /></el-icon>
                                <span :class="[isCollapse ? 'is-collapse' : '']">{{menu.name}}</span>
                            </template>
                            <el-menu-item class="aside-menu-childitem" v-for="child in menu.children" :key="child" :index="child.path">
                                <el-icon><component :is="child.icon" /></el-icon>
                                <template #title>
                                    {{child.name}}
                                </template>
                            </el-menu-item>
                        </el-sub-menu>
                    </div>
                </el-menu>
            </el-aside>
            <el-container>
                <el-header class="header" >
                    <el-row :gutter="20">
                        <el-col :span="1">
                            <div class="header-collapse" @click="onCollapse">
                                <el-icon><component :is="isCollapse ? 'expand':'fold'" /></el-icon>
                            </div>
                        </el-col>
                        <el-col :span="10" >
                            <div class="header-breadcrumb">
                                <el-breadcrumb separator="/" v-if="this.$route.matched[0].path != '/main'">
                                    <el-breadcrumb-item :to="{ path: '/' }">工作台</el-breadcrumb-item>
                                    <template v-for="(matched,m) in this.$route.matched" :key="m">
                                        <el-breadcrumb-item v-if="matched.name != undefined" >
                                        {{ matched.name }}
                                        </el-breadcrumb-item>
                                    </template>
                                </el-breadcrumb>
                                <el-breadcrumb separator="/" v-else>
                                    <el-breadcrumb-item>工作台</el-breadcrumb-item>
                                </el-breadcrumb> 
                            </div>
                        </el-col>
                        <el-col class="header-menu" :span="13">
                            <el-dropdown>
                                <div class="header-dropdown">
                                    <el-image class="avator-image" :src="avator" />
                                    <span>{{ username }}</span>
                                </div>
                                <template #dropdown>
                                    <el-dropdown-menu>
                                        <el-dropdown-item icon="el-icon-switch-button" @click="logout()">退出</el-dropdown-item>
                                        <el-dropdown-item icon="el-icon-unlock">修改密码</el-dropdown-item>
                                    </el-dropdown-menu>
                                </template>
                            </el-dropdown>
                        </el-col>
                    </el-row>
                </el-header>
                <el-main class="main">
                    <router-view></router-view>
                </el-main>
                <el-footer class="footer">
                    <el-icon style="width:2em;top:3px;font-size:18px"><place/></el-icon>
                    <a class="footer el-icon-place">2022 DevOps </a>
                </el-footer>
                <el-backtop target=".el-main"></el-backtop>
            </el-container>
        </el-container>
    </div>
</template>

<script>
import {useRouter} from 'vue-router'
export default {
    data() {
        return {
            avator: require('@/assets/avator/avator.png'),
            logo: require('@/assets/k8s/k8s-metrics.png'),
            isCollapse: false,
            asideWidth: '220px',
            routers: [],
        }
    },
    computed: {
        username() {
            let username = localStorage.getItem('username');
            return username ? username : '未知';
        },
    },
    methods: {
        onCollapse() {
            if (this.isCollapse) {
                this.asideWidth = '220px'
                this.isCollapse = false
            } else {
                this.isCollapse = true
                this.asideWidth = '64px'
            }
        },
        logout() {
            localStorage.removeItem('username');
            localStorage.removeItem('token');
            this.$router.push('/login');
        }
    },
    beforeMount() {
        this.routers = useRouter().options.routes
    }
}
</script>

<style scoped>
    .aside {
        transition: all .5s;
        background-color: #131b27;
    }
    .aside-logo {
        background-color: #131b27;
        height: 60px;
        color: white;
        cursor: pointer;
    }
    .logo-image {
        width: 40px;
        height: 40px;
        top: 12px;
        padding-left: 12px;
    }
    .logo-name {
        font-size: 20px;
        font-weight: bold;
        padding: 10px;
    }
    .aside::-webkit-scrollbar {
        display: none;
    }
    .aside-affix {
        border-bottom-width: 0;
    }
    .aside-menu {
        border-right-width: 0;
    }
    .aside-menu-item.is-active {
        background-color: #1f2a3a;
    }
    .aside-menu-item {
        padding-left: 20px !important;
    }
    .aside-menu-childitem {
        padding-left: 20px !important;
    }
    .aside-menu-childitem.is-active {
        background-color: #1f2a3a;
    }
    .aside-menu-childitem:hover {
        background-color: #142c4e;
    }
    .header {
        z-index:1200;
        line-height: 60px;
        font-size: 24px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, .12),0 0 6px rgba(0, 0, 0, .04)
    }
    .header-collapse {
        cursor: pointer;
    }
    .header-breadcrumb {
        padding-top: 0.9em;
    }
    .header-menu {
        text-align: right;
    }
    .is-collapse {
        display: none;
    }
    .header-dropdown {
        line-height: 60px;
        cursor: pointer;
    }
    .avator-image {
        top: 12px;
        width: 40px;
        height: 40px;
        border-radius: 50%;
        margin-right: 8px;
    }
    .main {
        padding: 10px;
    }
    .footer {
        z-index: 1200;
        color: rgb(187, 184, 184);
        font-size: 14px;
        text-align: center;
        line-height: 60px;
    }
</style>
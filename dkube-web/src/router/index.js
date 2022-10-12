import {createRouter,createWebHistory} from 'vue-router'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import Layout from "@/layout/Layout"

import jwt from 'jsonwebtoken'

const routes =[
    {
        path: '/login',
        component: () => import('@/views/login/Login.vue'),
        icon: "odometer",
        meta: {title: "登录", requireAuth: false},
    },
    {
        path: '/',
        redirect: '/home'
    },
    {
        path: '/home',
        component: Layout,
        icon: "odometer",
        children: [
            {
                path: "/home",
                name: "K8s概览",
                icon: "odometer",
                meta: {title: "K8s概览", requireAuth: true},
                component: () => import('@/views/home/Home.vue'),
            }
        ]
    },
    {
        path: '/workflow',
        component: Layout,
        icon: "VideoPlay",
        children: [
            {
                path: "/workflow",
                name: "工作流",
                icon: "VideoPlay",
                meta: {title: "工作流", requireAuth: true},
                component: () => import('@/views/workflow/Workflow.vue')
            }
        ]
    },
    {
        path: "/cluster",
        name: "集群信息",
        component: Layout,
        icon: "home-filled",
        meta: {title: "集群信息", requireAuth: true},
        children: [
            {
                path: "/cluster/node",
                name: "Node",
                icon: "el-icon-s-data",
                meta: {title: "Node", requireAuth: true},
                component: () => import("@/views/node/Node.vue")
            },
            {
                path: "/cluster/namespace",
                name: "Namespace",
                icon: "el-icon-document-add",
                meta: {title: "Namespace", requireAuth: true},
                component: () => import("@/views/namespace/Namespace.vue")
            },
            {
                path: "/cluster/persistentvolume",
                name: "PersistentVolume",
                icon: "el-icon-document-add",
                meta: {title: "PersistemtVolume", requireAuth: true},
                component: () => import("@/views/persistentvolume/PersistentVolume.vue")
            }
        ]
    },
    {
        path: '/workload',
        name: '工作负载',
        component: Layout,
        icon: "menu",
        meta: {title: "工作负载", requireAuth: true},
        children: [
            {
                path: '/workload/deployment',
                name: 'Deployment',
                icon: 'el-icon-s-data',
                meta: {title: "Deployment", requireAuth: true},
                component: () => import("@/views/deployment/Deployment.vue")
            },
            {
                path: '/workload/pod',
                name: 'Pod',
                icon: 'el-icon-document-add', 
                meta: {title: "Pod", requireAuth: true},
                component: () => import("@/views/pod/Pod.vue")
            },
            {
                path: '/workload/daemonset',
                name: 'DaemonSet',
                icon: 'el-icon-document-add',
                meta: {title: "DaemonSet", requireAuth: true},
                component: () => import("@/views/daemonset/DaemonSet.vue")
            },
            {
                path: '/workload/statefulset',
                name: 'StatefulSet',
                icon: 'el-icon-document-add', 
                meta: {title: "StatefulSet", requireAuth: true},
                component: () => import("@/views/statefulset/StatefulSet.vue")
            },
        ]
    },
    {
        path: "/loadbalance",
        name: "负载均衡",
        component: Layout,
        icon: "files",
        meta: {title: "负载均衡", requireAuth: true},
        children: [
            {
                path: "/loadbalance/service",
                name: "Service",
                icon: "el-icon-s-data",
                meta: {title: "Service", requireAuth: true},
                component: () => import("@/views/service/Service.vue")
            },
            {
                path: "/loadbalance/ingress",
                name: "Ingress",
                icon: "el-icon-document-add",
                meta: {title: "Ingress", requireAuth: true},
                component: () => import("@/views/ingress/Ingress.vue")
            }
        ]
    },
    {
        path: "/storage",
        name: "存储与配置",
        component: Layout,
        icon: "tickets",
        meta: {title: "存储与配置", requireAuth: true},
        children: [
            {
                path: "/storage/configmap",
                name: "Configmap",
                icon: "el-icon-document-add",
                meta: {title: "Configmap", requireAuth: true},
                component: () => import("@/views/configmap/ConfigMap.vue")
            },
            {
                path: "/storage/secret",
                name: "Secret",
                icon: "el-icon-document-add",
                meta: {title: "Secret", requireAuth: true},
                component: () => import("@/views/secret/Secret.vue")
            },
            {
                path: "/storage/persistentvolumeclaim",
                name: "PersistentVolumeClaim",
                icon: "el-icon-s-data",
                meta: {title: "PersistentVolumeClaim", requireAuth: true},
                component: () => import("@/views/persistentvolumeclaim/PersistentVolumeClaim.vue")
            },
        ]
    },
    {
        path: '/404',
        component: () => import('@/views/common/404.vue'),
        meta: {
            title: '404'
        }
    },
    {
        path: '/403',
        component: () => import('@/views/common/403.vue'),
        meta: {
            title: '403'
        }
    },
    {
        path: '/:pathMatch(.*)',
        redirect: '/404'
    },
]

const router = createRouter ({
    history: createWebHistory(),
    routes
})

NProgress.inc(100)
NProgress.configure({ easing: 'ease', speed: 600,showSpinner: false})

router.beforeEach((to,from,next) => {
    NProgress.start()
    if (to.meta.title) {
        document.title = to.meta.title
    } else {
        document.title = "DKube"
    }
    next()
})

router.beforeEach((to, from, next) => {
    jwt.verify(localStorage.getItem('token'), 'devops', function (err) {
        if (to.path === '/login') {
            next()
        } else if (err) {
            next('/login');
        } else {
            next();
        }
    });
});

router.afterEach (() => {
    NProgress.done()
})

export default router
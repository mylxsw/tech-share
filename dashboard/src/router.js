import Vue from 'vue';
import Router from 'vue-router';
import Shares from './views/Shares';
import Share from './views/Share';
import Login from './views/Login';
import Rank from './views/Rank';

Vue.use(Router);

const routerPush = Router.prototype.push;
Router.prototype.push = function push(location) {
    return routerPush.call(this, location).catch(error => error)
}

export default new Router({
    routes: [
        {path: '/', component: Shares},
        {path: '/share', component: Share},
        {path: '/login', component: Login},
        {path: '/rank', component: Rank},
    ]
});

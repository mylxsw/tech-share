import '@babel/polyfill'
import 'mutationobserver-shim'
import Vue from 'vue'
import './plugins/axios'
import './plugins/bootstrap-vue'
import App from './App.vue'
import router from './router'
import store from './store'
import uploader from 'vue-simple-uploader'

import { BootstrapVueIcons } from 'bootstrap-vue'

import { library } from '@fortawesome/fontawesome-svg-core'
import { faExternalLinkAlt, faPlus } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

import DateTime from "./components/DateTime";
import HumanTime from "./components/HumanTime";

library.add(faExternalLinkAlt);
library.add(faPlus)
Vue.component('font-awesome-icon', FontAwesomeIcon);

Vue.use(BootstrapVueIcons);
Vue.use(uploader);

Vue.component('DateTime', DateTime);
Vue.component('HumanTime', HumanTime);

Vue.config.productionTip = false;

const errorHandler = (error) => {
    if (error === 'access-denied') {
        window.location.href = '/#/login';
        return;
    }

    console.log(error);
}

Vue.config.errorHandler = errorHandler;
Vue.prototype.$throw = (error) => errorHandler(error, this);

Vue.prototype.QueryArgs = (route, name) => {
    return route.query[name] !== undefined ? route.query[name] : '';
}

/**
 * @return {string}
 */
Vue.prototype.ParseError = function (error) {
    if (error.response !== undefined) {
        if (error.response.status === 401) {
            this.$throw("access-denied");
        } 

        if (error.response.data !== undefined) {
            return error.response.data.error;
        }
    }

    return error.toString();
};

Vue.prototype.ToastSuccess = function (message) {
    this.$bvToast.toast(message, {
        title: 'OK',
        variant: 'success',
        autoHideDelay: 3000,
        toaster: 'b-toaster-top-center',
    });
};

Vue.prototype.ToastError = function (message) {
    this.$bvToast.toast(this.ParseError(message), {
        title: 'ERROR',
        variant: 'danger',
        autoHideDelay: 3000,
        toaster: 'b-toaster-top-center',
    });
};

Vue.prototype.SuccessBox = function (message, cb) {
    cb = cb || function () {};
    this.$bvModal.msgBoxOk(message, {
        title: '操作成功',
        centered: true,
        okVariant: 'success',
        headerClass: 'p-2 border-bottom-0',
        footerClass: 'p-2 border-top-0',
    }).then(cb);
};

Vue.prototype.ErrorBox = function (message, cb) {
    cb = cb || function () {};
    this.$bvModal.msgBoxOk(this.ParseError(message), {
        centered: true,
        title:'出错了',
        okVariant: 'danger',
        headerClass: 'p-2 border-bottom-0',
        footerClass: 'p-2 border-top-0',
    }).then(cb);
};

new Vue({
    router,
    store,
    render: h => h(App)
}).$mount('#app');

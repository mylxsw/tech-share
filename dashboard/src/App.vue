<template>
  <div id="app">
    <b-container fluid>
        <b-navbar type="dark" toggleable="md" variant="dark" class="mb-3" sticky>
            <b-navbar-brand href="/">Tech Share <a href="https://github.com/mylxsw/tech-share" class="text-white" style="font-size: 30%">{{ version }}</a></b-navbar-brand>
            <b-collapse is-nav id="nav_dropdown_collapse">
                <ul class="navbar-nav flex-row ml-md-auto d-none d-md-flex"></ul>
                <b-navbar-nav>
                    <b-nav-item href="/" exact v-if="isLogon()">分享列表</b-nav-item>
                    <b-nav-item @click="logout()" v-if="isLogon()">退出</b-nav-item>
                </b-navbar-nav>
            </b-collapse>
        </b-navbar>
        <div class="main-view">
            <router-view/>
        </div>
    </b-container>
    
  </div>
</template>

<script>
    import axios from 'axios';

    export default {
        data() {
            return {
                version: 'v-0',
            }
        },
        methods: {
            isLogon() {
                return this.$store.getters.user != null && this.$store.getters.user.id > 0;
            },
            logout() {
                axios.post('/api/auth/logout').then(response => {
                    this.$store.commit('removeUser');
                    window.location.reload();
                }).catch(err => {this.ErrorBox(err)});
            }
        },
        mounted() {
            axios.get('/api/inspect/version').then(response => {
                this.version = response.data.version;
            });
        },
        beforeMount() {
            axios.get('/api/auth/current').then(response => {
                this.$store.commit('updateUser', response.data);
            }).catch(err => {this.ToastError(err)})
        }
    }
</script>

<style>
    .container-fluid {
        padding: 0;
    }

    .main-view {
        padding: 15px;
    }

    .th-column-width-limit {
        max-width: 300px;
    }

    @media screen and (max-width: 1366px) {
        .th-autohide-md {
            display: none;
        }
    }
    @media screen and (max-width: 768px) {
        .th-autohide-sm {
            display: none;
        }
        .search-box {
            display: none;
        }
    }

</style>

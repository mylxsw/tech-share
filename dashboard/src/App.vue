<template>
  <div id="app">
    <b-container fluid>
        <b-navbar type="dark" toggleable="md" variant="dark" class="mb-3" sticky>
            <b-navbar-brand href="/">Tech Share <a href="https://github.com/mylxsw/tech-share" class="text-white" style="font-size: 30%">{{ version }}</a></b-navbar-brand>
            <b-collapse is-nav id="nav_dropdown_collapse">
                <ul class="navbar-nav flex-row ml-md-auto d-none d-md-flex"></ul>
                <b-navbar-nav>
                    <b-nav-item href="/#/?act=" exact v-if="isLogon()">所有分享</b-nav-item>
                    <b-nav-item href="/#/?act=recently" exact v-if="isLogon()">最新分享</b-nav-item>
                    <b-nav-item href="/#/?act=my" exact v-if="isLogon()">我发起的</b-nav-item>
                    <b-nav-item href="/#/rank" exact v-if="isLogon()">排行榜</b-nav-item>
                    <b-nav-item-dropdown right v-if="isLogon()">
                        <template #button-content>
                            <em>{{ $store.getters.user.name }}</em>
                        </template>
                        <b-dropdown-item @click="logout()">退出</b-dropdown-item>
                    </b-nav-item-dropdown>
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
                axios.post('/api/auth/logout').then(() => {
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
            }).catch(() => {})
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

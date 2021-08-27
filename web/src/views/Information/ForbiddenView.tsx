import React, { useEffect, Fragment, ReactNode } from "react";

import { Grid, makeStyles } from "@material-ui/core";

import { useNotifications } from "@hooks/NotificationsContext";
import { useRedirectionURL } from "@hooks/RedirectionURL";
import { useUserPreferences as userUserInfo } from "@hooks/UserInfo";
import InformationLayout from "@layouts/InformationLayout";
import LoadingPage from "@views/LoadingPage/LoadingPage";

export interface Props {
    title: string;
}

const ForbiddenView = function (props: Props) {
    const classes = useStyles();
    const { createErrorNotification, resetNotification } = useNotifications();
    const [resp, fetch, , err] = userUserInfo();
    const redirectionURL = useRedirectionURL();

    useEffect(() => {
        if (err) {
            createErrorNotification("Error");
            console.error(`Error: ${err.message}`);
        }
    }, [resetNotification, createErrorNotification, err]);

    useEffect(() => {
        fetch();
    }, [fetch]);

    return (
        <ComponentOrLoading ready={resp !== undefined}>
            <InformationLayout id="status" title={props.title} showBrand>
                <Grid container className={classes.container}>
                    <Grid item xs={12}>
                        <div style={{ textAlign: "left" }}>
                            You are forbidden access to the site <a href={redirectionURL}>{redirectionURL}</a>. Please
                            contact an Administrator if you think this is in error {resp?.display_name}.
                        </div>
                    </Grid>
                </Grid>
            </InformationLayout>
        </ComponentOrLoading>
    );
};

const useStyles = makeStyles((theme) => ({
    container: {
        paddingTop: theme.spacing(4),
        paddingBottom: theme.spacing(4),
        display: "block",
        justifyContent: "center",
    },
}));

export default ForbiddenView;

interface ComponentOrLoadingProps {
    ready: boolean;

    children: ReactNode;
}

function ComponentOrLoading(props: ComponentOrLoadingProps) {
    return (
        <Fragment>
            <div className={props.ready ? "hidden" : ""}>
                <LoadingPage />
            </div>
            {props.ready ? props.children : null}
        </Fragment>
    );
}
